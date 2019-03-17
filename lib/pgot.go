package pgot

import (
	"bytes"
	"encoding/json"
	"errors"
	"git.lenzplace.org/lenzj/chunkio"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	templateName = "."
	includeKey   = "pgotInclude"
)

var (
	key                     = []byte(";;;\n")
	ErrReservedTemplateName = errors.New("attempt to include reserved filename \"" + templateName + "\"")
	ErrMalformedFrontMatter = errors.New("frontmatter is missing or malformed")
	ErrMalformedInclude     = errors.New(includeKey + " is malformed")
)

type Parser struct {
	gin []string           // Colon separated string containing paths to search
	fms *fmStack           // Frontmatter stack
	t   *template.Template // Golang template
	td  []byte             // Byte array containing main template to parse
}

// parses frontmatter from byte array into frontmatter object
// data : a byte array containing frontmatter (json only, no keys)
// fm : the frontMatter is stored into this variable
func parseFMData(data []byte, fm *frontMatter) error {
	data = bytes.TrimSpace(data)

	if len(data) == 0 {
		return nil
	}

	var jdata interface{}
	err := json.Unmarshal(data, &jdata)
	if err != nil {
		return err
	}

	var fmd map[string]interface{}
	if jdata == nil {
		fmd = map[string]interface{}{}
	} else {
		fmd = jdata.(map[string]interface{})
	}

	if fm.ns == "." {
		//Merge with root namespace
		for k, v := range fmd {
			fm.fm[k] = v
		}
	} else {
		fm.fm[fm.ns] = fmd
	}
	return nil
}

// parses a new frontmatter string into the Parser object
// fm : a string containing frontmatter (json only, no keys)
// fd : directory path where stream is being read from
// ns : namespace to import frontmatter into
func (c *Parser) ParseFMString(fmstr, fd, ns string) error {
	if err := parseFMData([]byte(fmstr), c.fms.first()); err != nil {
		return err
	}
	return nil
}

// parses a new frontmatter stream into the Parser object
// cio : input stream to read from
// fd : directory path where stream is being read from
// ns : namespace to import frontmatter into
func (c *Parser) ParseFMStream(cio *chunkio.Reader, fd, ns string) error {
	cio.SetKey(key)
	// There shouldn't be any data prior to the first frontmatter key
	fmstring, err := ioutil.ReadAll(cio)
	if len(fmstring) > 0 {
		return ErrMalformedFrontMatter
	}
	// Reset stream and read frontmatter
	cio.Reset()
	fmstring, err = ioutil.ReadAll(cio)
	if err != nil {
		return err
	}
	c.fms.push(&frontMatter{
		fd: fd,
		fm: map[string]interface{}{},
		ns: ns,
	})
	if err = parseFMData([]byte(fmstring), c.fms.last()); err != nil {
		return err
	}
	return nil
}

// Creates and returns a new Parser object
// in : input stream to read from
// fd : default base directory path where stream is being read from
// gin : alternate directory paths to try when including relative paths
func NewParser(in io.Reader, fd string, gin []string) (*Parser, error) {
	cio := chunkio.NewReader(in)
	fmstack := newfmStack()
	c := Parser{
		gin: gin,
		fms: fmstack,
		t:   template.New(templateName).Funcs(funcMap),
		td:  nil,
	}
	err := c.ParseFMStream(cio, fd, ".")
	if err != nil {
		return nil, err
	}
	cio.Reset()
	cio.SetKey(nil)
	c.td, _ = ioutil.ReadAll(cio)
	return &c, nil
}

// If a pgotInclude array exists in frontmatter then this function will process
// it.
func (c *Parser) ProcessFMInclude() error {
	fm := c.fms.last()
	if cgIncludes, ok := fm.fm[includeKey]; ok {
		delete(fm.fm, includeKey)
		if cgi, ok := cgIncludes.([]interface{}); ok {
			for _, in := range cgi {
				switch item := in.(type) {
				case string:
					if err := c.IncludeFile(item, fm.fd, "."); err != nil {
						return err
					}
				case map[string]interface{}:
					if err := c.IncludeFile(item["f"].(string), fm.fd, item["n"].(string)); err != nil {
						return err
					}
				default:
					return ErrMalformedInclude
				}
			}
		} else {
			return ErrMalformedInclude
		}
	}
	return nil
}

// Include a got file.  This will process any frontmatter in the file and will
// also create a template with the same name as the file.
// f: path to file for import
// pp: path of parent initiating the import
// ns: namespace to import frontmatter into
func (c *Parser) IncludeFile(f, pp, ns string) error {
	base := filepath.Base(f)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	if name == templateName {
		return ErrReservedTemplateName
	}

	var (
		fp  string
		in  *os.File
		err error
	)

	if filepath.IsAbs(f) {
		if in, err = os.Open(f); err != nil {
			return err
		}
		fp = f
		defer in.Close()
	} else {
		var (
			inFile  string
			relFile string
			pwd     string
		)
		for _, p := range c.gin {
			if p == "" {
				inFile = filepath.Join(pp, f)
			} else {
				inFile = filepath.Join(p, f)
			}
			pwd, err = os.Getwd()
			if err != nil {
				return err
			}
			relFile, err = filepath.Rel(pwd, inFile)
			if err != nil {
				relFile = inFile
			}
			if in, err = os.Open(relFile); os.IsNotExist(err) {
				continue
			} else {
				fp = inFile
				defer in.Close()
				break
			}
		}
		if err != nil {
			return err
		}
	}
	cio := chunkio.NewReader(in)
	if err = c.ParseFMStream(cio, filepath.Dir(fp), ns); err != nil {
		return err
	}
	if err = c.ProcessFMInclude(); err != nil {
		return err
	}
	if err = c.parseTemplate(cio, name); err != nil {
		return err
	}
	return nil
}

// This function parses template content into the specified template name.
// cio : This is a chunkio input stream to parse from.  This function assumes
//       that frontmatter has already been read from this stream.
// tn  : Content is parsed into a template using namespace tn.
func (c *Parser) parseTemplate(cio *chunkio.Reader, tn string) error {
	// This assumes the frontmatter section has already completed parsing
	cio.Reset()
	cio.SetKey(nil)
	tdata, err := ioutil.ReadAll(cio)
	if err != nil && err != io.ErrUnexpectedEOF {
		return err
	}
	_, err = c.t.New(tn).Parse(string(tdata))
	if err != nil {
		return err
	}
	return nil
}

// This function executes the template and sends the output to "out".
func (c *Parser) Execute(out io.Writer) error {
	// Flatten the frontmatter stack
	fm_flat := map[string]interface{}{}
	for item := c.fms.pop(); item != nil; item = c.fms.pop() {
		for k, v := range item.fm {
			fm_flat[k] = v
		}
	}
	c.fms.push(&frontMatter{fm: fm_flat})

	if _, err := c.t.Parse(string(c.td)); err != nil {
		return err
	}
	if err := c.t.Execute(out, c.fms.last().fm); err != nil {
		return err
	}
	return nil
}
