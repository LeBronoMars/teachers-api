package models

type QryClassSubjects struct {
	ClassSubjectId string `json:"class_subject_id"`
	ClassSubjectIsSynced bool `json:"class_subject_id_is_synced"`
	ClassSubjectCreatedBy string `json:"class_subject_created_by"`
	ClassSubjectDeletedAt string `json:"class_subject_deleted_at"`
	ClassId string `json:"class_id"`
	ClassIsSynced string `json:"class_is_synced"`
	ClassGradeLevel string `json:"class_grade_level"`
	ClassSection string `json:"class_section"`
	ClassSchoolYearFrom int `json:"class_school_year_from"`
	ClassSchoolYearTo int `json:"class_school_year_to"`
	ClassSchool string `json:"class_school"`
	ClassRemarks string `json:"class_remarks"`
	TeacherIsSynced bool `json:"teacher_is_synched"`
	TeacherEmployeeNo string `json:"teacher_employee_no"`
	TeacherFirstName string `json:"teacher_first_name"`
	TeacherMiddleName string `json:"teacher_middle_name"`
	TeacherLastName string `json:"teacher_last_name"`
	TeacherBirthDate string `json:"teacher_birth_date"`
	TeacherBirthPlace string `json:"teacher_birth_place"`
	TeacherGender string `json:"teacher_gender"`
	TeacherCivilStatus string `json:"teacher_civil_status"`
	TeacherEmail string `json:"teacher_email"`
	TeacherAddress string `json:"teacher_address"`
	TeacherContactNo string `json:"teacher_contact_no"`
	TeacherStatus string `json:"teacher_status"`
	TeacherUserRole string `json:"teacher_user_role"`
	TeacherPosition string `json:"teacher_position"`
	TeacherPicUrl string `json:"teacher_pic_url"`
	SubjectId string `json:"subject_id"`
	SubjectIsSynced bool `json:"subject_is_synced"`
	SubjectName string `json:"subject_name"`
	SubjectCode string `json:"subject_code"`
	SubjectDescription string `json:"subject_description"`
}