package main

import(
	"os"
	"flag"
	"fmt"
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

func main() {
	if path == "" {
		fmt.Println("Error: Base file path undefined!")
		return
	}

	xmlFile, err := os.Open(path + "/profiles/Accounting.profile")
	if err != nil {
		fmt.Println("Error opening XML file: ", err)
		return
	}
	defer xmlFile.Close()

	b, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		fmt.Println("Error reading XML into byte slice: ", err)
		return
	}

	//var p Profile
	p := new(Profile)
	xml.Unmarshal(b, p)

	fmt.Println("Profile License Type: ", p.UserLicense)
	for i, f := range p.FieldPermList {
		fmt.Printf("\tOld Value - Index: %d - Field: %v\n", i, f.Field)
		if !f.Editable {
			f.Editable = true
			fmt.Println("\t\tSetting editable field permission")
		} 
		if !f.Readable {
			f.Readable = true
			fmt.Println("\t\tSetting readable field permission")
		}
	}
	
}