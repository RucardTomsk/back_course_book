package model

type FullPlan struct {
	Guid string `json:"guid"`
	//Рабочая программа дисциплины
	Name string `json:"name"`
	//Направлению подготовки
	DirectionTraining string `json:"directionTraining"`
	//(профиль) подготовки
	TrainingProfile string `json:"trainingProfile"`
	//Форма обучения
	FormTraining string `json:"formTraining"`
	//Квалификация
	Qualification string `json:"qualification"`
	//Год приема
	Year string `json:"year"`
	//Код дисциплины
	Code string `json:"code"`
	//Целью освоения дисциплины
	PurposeMastering string `json:"purposeMastering"`
	//Результатами освоения дисциплины
	ResultsMastering string `json:"resultsMastering"`
	//Место дисциплины
	PlaceDiscipline string `json:"placeDiscipline"`
	//Семестр(ы) освоения
	SemesterMastering string `json:"semesterMastering"`
	//Язык реализации
	ImplementationLanguage string `json:"implementationLanguage"`
	//Задачи освоения дисциплины
	DevelopmentMastering string `json:"developmentMastering"`
	//Входные требования для освоения дисциплины
	EntranceRequirements string `json:"entranceRequirements"`
	//Объем дисциплины (модуля)
	ScopeDiscipline string `json:"scopeDiscipline"`
	//Содержание дисциплины
	ContentDiscipline string `json:"contentDiscipline"`
	//Текущий контроль по дисциплине
	CurrentControl string `json:"currentControl"`
	//Порядок проведения и критерии оценивания промежуточной аттестации
	EvaluationProcedure string `json:"evaluationProcedure"`
	//Учебно-методическое обеспечение
	MethodologicalSupport string `json:"methodologicalSupport"`
	//Перечень учебной литературы и ресурсов сети Интернет
	References string `json:"references"`
	//Перечень информационных технологий
	ListInformationTechnologies string `json:"listInformationTechnologies"`
	//Материально-техническое обеспечение
	MaterialSupport string `json:"materialSupport"`
	//Информация о разработчиках
	InformationDevelopers string `json:"informationDevelopers"`
	//Информация о разработчиках
	NameFaculty string `json:"nameFaculty"`
}

type BriefPlan struct {
	//Код дисциплины
	Code string `json:"code"`
	//Название дисциплины
	Name string `json:"name"`
	//guid
	Guid string `json:"guid"`
	//Семетры
	//Семестр(ы) освоения
	SemesterMastering string
}
