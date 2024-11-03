package main

import (
	"encoding/json"
	"fmt"
	"net"
)

const (
	CreateSchoolMethod      = "/school/create"
	CreateClassMethod       = "/class/create"
	CreatePersonMethod      = "/person/create"
	AddStudentToClassMethod = "/class/add/student"
	WhoAmIMethod            = "/who/am/i"
)

type Server interface {
	Start(port string) error
	Stop() error
}

type server struct {
	listener        net.Listener
	schools         map[uint]School
	people          map[uint]Person
	classes         map[uint]Class
	schoolIdCounter uint
	personIdCounter uint
	classIdCounter  uint
}

func NewServer() Server {
	return &server{
		schools:         make(map[uint]School),
		people:          make(map[uint]Person),
		classes:         make(map[uint]Class),
		schoolIdCounter: 1,
		personIdCounter: 1,
		classIdCounter:  1,
	}
}

func (s *server) Start(port string) error {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	s.listener = ln
	fmt.Println("Server started on port", port)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
				fmt.Println("Server stopped")
				return nil
			}
			return err
		}
		go s.handleConnection(conn)
	}
}

func (s *server) Stop() error {
	if s.listener != nil {
		fmt.Println("Stopping server...")
		return s.listener.Close()
	}
	return nil
}

func (s *server) handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		var req Request
		if err := decoder.Decode(&req); err != nil {
			fmt.Println("Failed to decode request:", err)
			return
		}

		var resp Response
		switch req.Method {
		case CreateSchoolMethod:
			resp = s.createSchool(req.Data)
		case CreateClassMethod:
			resp = s.createClass(req.Data)
		case CreatePersonMethod:
			resp = s.createPerson(req.Data)
		case AddStudentToClassMethod:
			resp = s.addStudentToClass(req.Data)
		case WhoAmIMethod:
			resp = s.whoAmI(req.Data)
		default:
			resp = Response{Status: false, Message: "Invalid method"}
		}

		if err := encoder.Encode(&resp); err != nil {
			fmt.Println("Failed to encode response:", err)
			return
		}
	}
}

func (s *server) createSchool(data interface{}) Response {
	var school School
	if !s.decodeRequest(data, &school) {
		return Response{Status: false, Message: "Invalid school data"}
	}

	school.Id = s.schoolIdCounter
	s.schoolIdCounter++
	s.schools[school.Id] = school

	return Response{Status: true, Message: "School created successfully", Data: school}
}

func (s *server) createPerson(data interface{}) Response {
	var person Person
	if !s.decodeRequest(data, &person) {
		return Response{Status: false, Message: "Invalid person data"}
	}

	person.Id = s.personIdCounter
	s.personIdCounter++
	s.people[person.Id] = person

	return Response{Status: true, Message: "Person created successfully", Data: person}
}

func (s *server) createClass(data interface{}) Response {
	var class Class
	if !s.decodeRequest(data, &class) {
		return Response{Status: false, Message: "Invalid class data"}
	}

	if _, exists := s.schools[class.SchoolId]; !exists {
		return Response{Status: false, Message: "School not found"}
	}

	teacher, exists := s.people[class.Teacher.Id]
	if !exists || s.teacherIsStudent(teacher.Id) {
		return Response{Status: false, Message: "Invalid or conflicting teacher"}
	}

	class.Id = s.classIdCounter
	s.classIdCounter++
	teacher.Classes = append(teacher.Classes, class.Id)
	class.Teacher = teacher

	s.updateTeacherClasses(teacher)
	s.people[teacher.Id] = teacher
	s.classes[class.Id] = class

	school := s.schools[class.SchoolId]
	school.Classes = append(school.Classes, class)
	s.schools[class.SchoolId] = school

	class.Teacher.Classes = nil
	return Response{Status: true, Message: "Class created successfully", Data: class}
}

func (s *server) addStudentToClass(data interface{}) Response {
	var req AddStudentToClassReq
	if !s.decodeRequest(data, &req) {
		return Response{Status: false, Message: "Invalid request data"}
	}

	class, classExists := s.classes[req.ClassId]
	if !classExists {
		return Response{Status: false, Message: "Class not found"}
	}

	student, studentExists := s.people[req.StudentId]
	if !studentExists || s.isInOtherSchoolClasses(student.Id, class.SchoolId) {
		return Response{Status: false, Message: "Student conflict or not found"}
	}

	if s.isAlreadyInClass(student.Id, class.Id) {
		return Response{Status: false, Message: "Student already in this class"}
	}

	student.Classes = append(student.Classes, class.Id)
	class.Students = append(class.Students, student)

	s.people[student.Id] = student
	s.classes[class.Id] = class
	s.updateClassInSchool(class)

	return Response{Status: true, Message: "Student added to class", Data: student}
}

func (s *server) whoAmI(data interface{}) Response {
	var person Person
	if !s.decodeRequest(data, &person) {
		return Response{Status: false, Message: "Invalid person data"}
	}

	actualPerson, exists := s.people[person.Id]
	if !exists {
		return Response{Status: false, Message: "Person not found"}
	}

	return Response{Status: true, Message: "Person found", Data: actualPerson}
}

func (s *server) teacherIsStudent(teacherId uint) bool {
	for _, class := range s.classes {
		for _, student := range class.Students {
			if student.Id == teacherId {
				return true
			}
		}
	}
	return false
}

func (s *server) isInOtherSchoolClasses(studentId, schoolId uint) bool {
	for _, class := range s.classes {
		if class.SchoolId != schoolId {
			for _, student := range class.Students {
				if student.Id == studentId {
					return true
				}
			}
		}
	}
	return false
}

func (s *server) isAlreadyInClass(studentId, classId uint) bool {
	for _, studentClasses := range s.people[studentId].Classes {
		if studentClasses == classId {
			return true
		}
	}
	return false
}

func (s *server) updateTeacherClasses(teacher Person) {
	for id, existingClass := range s.classes {
		if existingClass.Teacher.Id == teacher.Id {
			existingClass.Teacher.Classes = teacher.Classes
			s.classes[id] = existingClass
		}
	}
}

func (s *server) updateClassInSchool(class Class) {
	school := s.schools[class.SchoolId]
	for i := range school.Classes {
		if school.Classes[i].Id == class.Id {
			school.Classes[i] = class
			break
		}
	}
	s.schools[class.SchoolId] = school
}

func (s *server) decodeRequest(data, v interface{}) bool {
	reqData, err := json.Marshal(data)
	if err != nil {
		return false
	}
	return json.Unmarshal(reqData, v) == nil
}
