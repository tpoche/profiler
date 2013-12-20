package main

import(
	"os"
	"errors"
	"flag"
	"fmt"
	"strings"
	"encoding/xml"
	"io/ioutil"
)

type Profile struct {
	FieldPermList  []FieldPermissions 		`xml:"fieldPermissions"`
	ObjectPermList []ObjectPermissions		`xml:"objectPermissions"`
	RecordTypeList []RecordTypeVisibilities	`xml:"recordTypeVisibilities"`
	UserLicense     string					`xml:"userLicense"`
}

type FieldPermissions struct {
	Editable bool	`xml:"editable"`
	Field 	 string	`xml:"field"`
	Readable bool	`xml:"readable"`
}

type ObjectPermissions struct {
	AllowCreate 	 bool   `xml:"allowCreate"`
	AllowDelete 	 bool 	`xml:"allowDelete"`
	AllowEdit 		 bool	`xml:"allowEdit"`
	AllowRead 		 bool   `xml:"allowRead"`
	ModifyAllRecords bool   `xml:"modifyAllRecords"`
	Object 			 string	`xml:"object"`
	ViewAllRecords 	 bool	`xml:"viewAllRecords"`
}

type RecordTypeVisibilities struct {
	Default    bool   `xml:"default"`	
	RecordType string `xml:"recordType"`
	Visible    bool	  `xml:"visible"`
}

var path string

func init() {
	flag.StringVar(&path, "filepath", ".", "base filepath for the program")
}

func NewProfileFromFile(filepath string) (*Profile, error) {
	if filepath == "" {
		fmt.Println("Error: Base file path undefined!")
		return nil, errors.New("NewProfileFromFile: Invalid path")
	}

	xmlFile, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening XML file: ", err)
		return nil, errors.New("NewProfileFromFile: Error opening file")
	}
	defer xmlFile.Close()

	b, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		fmt.Println("Error reading XML into byte slice: ", err)
		return nil, errors.New("NewProfileFromFile: Error reading file into byte slice")
	}

	//var p Profile
	pro := new(Profile)
	xml.Unmarshal(b, pro)
	if err != nil {
		return nil, errors.New("NewProfileFromFile: Error unmarshaling xml")
	}
	
	return pro, nil
}

func (p *Profile) WriteToFile(filepath string) (int, error) {
	// check that path is valid
	folders := strings.Split(filepath, "/")
	basepath := strings.Join(folders[:len(folders)-1], "/")
	_, e := os.Stat(basepath)
	if os.IsNotExist(e) {
		return 0, errors.New("WriteProfileToFile: Invalid file path")
	}
	// write modified profile to file
	out, e := xml.MarshalIndent(p, "", "    ")
	if e != nil {
		fmt.Printf("Error marshaling XML: %v\n", e)
		return 0, e
	}

	fout, e := os.Create(filepath)
	if e != nil {
		fmt.Println("Error creating file: ", e)
		return 0, e
	}

	nb, e := fout.Write(out)
	if e != nil {
		fmt.Println("Error writing file: ", e)
		return 0, e
	} 

	return nb, nil
}

func (p *Profile) UpdateFieldPerms(objName string, rd bool, ed bool) (int, error) {
	if objName == "" {
		return 0, errors.New("UpdateFieldPerms: Object name must be specified")
	}
	count := 0
	for i, v := range p.FieldPermList {
		if strings.HasPrefix(v.Field, objName) {
			if v.Readable != rd {
				p.FieldPermList[i].Readable = rd
				count += 1
			}
			if v.Editable != ed {
				p.FieldPermList[i].Editable = ed
				count += 1
			}
		}		
	}
	return count, nil
}

func main() {	
	flag.Parse()
	p, err := NewProfileFromFile(path + "/profiles/Accounting.profile")
	if err != nil {
		fmt.Println(err)
		return
	}

	// process field permissions
	for i, f := range p.FieldPermList {
		fmt.Printf("\tOld Value - Index: %d - Field: %v\n", i, f.Field)
		if !f.Editable {
			//f.Editable = true
			p.FieldPermList[i].Editable = true
		} 
		if !f.Readable {
			//f.Readable = true
			p.FieldPermList[i].Readable = true
		}
	}

	n, err := p.WriteToFile(path + "/out/Accounting2.profile")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Modified profile file created with size: %d\n", n)	
}