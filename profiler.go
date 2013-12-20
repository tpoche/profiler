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
var objects string
var fileOutput bool

func init() {
	flag.StringVar(&path, "filepath", ".", "base filepath for the program")
	flag.StringVar(&objects, "o", "", "comma separated list of objects to perform processing on (defaults to all objects found in profile)")
	flag.BoolVar(&fileOutput, "f", false, "enable writing output to file")
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

func (p *Profile) UpdateFieldPerms(rd bool, ed bool) (int, error) {
	obs, err := p.mapUserObjects(objects)
	if obs == nil || err != nil {
		return 0, errors.New("failed to create object map")
	}
	count := 0
	for i, v := range p.FieldPermList {
		fullName := strings.Split(v.Field, ".")
		if len(fullName) != 2 {
			return 0, errors.New("invalid field permission name encountered")
		}
		objName  := fullName[0]
		if obs[objName] == true {
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

func (p *Profile) mapUserObjects(userObjs string) (map[string]bool, error) {
	objMap := make(map[string]bool)
	objsFromProfile, err := p.GetObjectsWithFieldPerms()
	if err != nil {
		return nil, err
	}
	if objects == "" {
		// all objects specified
		for _, v := range objsFromProfile {
			objMap[v] = true
		}
		return objMap, nil
	}
	
	objs := strings.Split(objects, ",")
	for _, v := range objs {
		for _, x := range objsFromProfile {
			if v == x {
				objMap[x] = true
			}
		}
	}
	return objMap, nil
}

func (p *Profile) GetObjectsWithFieldPerms() ([]string, error) {
	objMap := make(map[string]int)

	for _, v := range p.FieldPermList {
		fieldFull := strings.Split(v.Field, ".")
		if len(fieldFull) != 2 {
			return nil, errors.New("ListObjectsWithFieldPerms: Invalid field name")
		} 
		
		curr := fieldFull[0]
		objMap[curr] += 1
	}

	keys := make([]string, len(objMap))
	for k, _ := range objMap {
		keys = append(keys, k)
	}
	return keys, nil
}

func main() {	
	flag.Parse()
	p, err := NewProfileFromFile(path + "/profiles/BillingSupport.profile")
	if err != nil {
		fmt.Println(err)
		return
	}

	objs, err := p.GetObjectsWithFieldPerms()
	if err != nil {
		fmt.Println(err)
		return
	} 
	fmt.Println("Objects with Field Permissions found: ", objs)
	
	c, err := p.UpdateFieldPerms(true, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Updated field permissions successfully! Modified perm count: ", c)

	
	if fileOutput {
		n, err := p.WriteToFile(path + "/out/BillingSupport.profile")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Modified profile file created with size: %d\n", n)	
	}	
}