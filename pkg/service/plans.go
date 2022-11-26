package service

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/RucardTomsk/course_book/model"
	"github.com/RucardTomsk/course_book/pkg/repository"
	"github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx/v3"
)

const (
	key_direction_training          = "DirectionTraining"
	key_training_profile            = "TrainingProfile"
	key_form_training               = "FormTraining"
	key_qualification               = "Qualification"
	key_year_admission              = "Year"
	key_name                        = "Name"
	key_purpose_mastering           = "PurposeMastering"
	key_results_mastering           = "ResultsMastering"
	key_place_of_discipline         = "PlaceDiscipline"
	key_semester_development        = "SemesterMastering"
	key_ZE                          = "ZE"
	key_K                           = "K"
	key_L                           = "L"
	key_C                           = "C"
	key_P                           = "P"
	key_LR                          = "LR"
	key_code                        = "Code"
	key_ImplementationLanguage      = "ImplementationLanguage"
	key_DevelopmentMastering        = "DevelopmentMastering"
	key_EntranceRequirements        = "EntranceRequirements"
	key_ContentDiscipline           = "ContentDiscipline"
	key_CurrentControl              = "CurrentControl"
	key_EvaluationProcedure         = "EvaluationProcedure"
	key_MethodologicalSupport       = "MethodologicalSupport"
	key_References                  = "References"
	key_ListInformationTechnologies = "ListInformationTechnologies"
	key_MaterialSupport             = "MaterialSupport"
	key_InformationDevelopers       = "InformationDevelopers"
	key_NameFacylte                 = "NameFaculty"
)

type PlansService struct {
	repo repository.Plans
}

func NewPlansService(repo repository.Plans) *PlansService {
	return &PlansService{repo: repo}
}

type dictCompetenciesStruct struct {
	title    string
	dictComp map[string]string
}

func (s *PlansService) GetPlans(guid_programm string) (map[string][]model.BriefPlan, error) {
	mas_plan, err := s.repo.GetMasPlan(guid_programm)
	if err != nil {
		return nil, err
	}
	dict_sort_plan := make(map[string][]model.BriefPlan)
	test_key_mas := ""
	for _, value := range mas_plan {
		semestr_start := value.SemesterMastering
		concrent_semestr_mas := strings.Split(semestr_start, "\n")
		var mas_int_semestr []string
		for _, str := range concrent_semestr_mas {
			con_semestr := strings.Split(strings.Split(str, ",")[0], " ")[1]
			mas_int_semestr = append(mas_int_semestr, con_semestr)
		}
		var key_semester string
		if len(mas_int_semestr) > 1 {
			key_semester = fmt.Sprintf("Семестры %s-%s", mas_int_semestr[0], mas_int_semestr[len(mas_int_semestr)-1])
		} else {
			key_semester = fmt.Sprintf("Семестр %s", mas_int_semestr[0])
		}
		if !strings.Contains(test_key_mas, key_semester) {
			mas_brief := []model.BriefPlan{value}
			dict_sort_plan[key_semester] = mas_brief
			test_key_mas += key_semester + " "
		} else {
			dict_sort_plan[key_semester] = append(dict_sort_plan[key_semester], value)
		}
	}

	keys := make([]string, 0, len(dict_sort_plan))
	for k := range dict_sort_plan {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	end_dict := make(map[string][]model.BriefPlan)

	for _, k := range keys {
		end_dict[k] = dict_sort_plan[k]
	}

	return end_dict, nil
}

func (s *PlansService) GetWorkProgram(guid_plan string) (model.FullPlan, error) {
	return s.repo.GetWorkProgram(guid_plan)
}

func (s *PlansService) SavePlan(guid_plan string, key_field string, text string) error {
	return s.repo.SavePlan(guid_plan, key_field, text)
}

func (s *PlansService) GetField(guid_plan string, key_field string) (string, error) {
	return s.repo.GetField(guid_plan, key_field)
}

func (s *PlansService) GetNamePlans(guid string) ([]string, error) {
	return s.repo.GetNamePlans(guid)
}

func IsEngByLoop(str string) bool {
	for i := 0; i < len(str); i++ {
		if str[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func (s *PlansService) CreatePlans(NameDiscipline string, ByteTable []byte, guid_faculty string) error {
	logrus.Info("Start CreatePlans")

	wb, err := xlsx.OpenBinary(ByteTable)
	if err != nil {
		return err
	}
	fmt.Println("dict_competencies")
	dict_competencies, err := get_dict_competencies(wb)
	if err != nil {
		return nil
	}
	fmt.Println("sheetTitel")
	sheetTitel, ok := wb.Sheet["Титул"]
	if !ok {
		return errors.New("Sheet not found")
	}

	CellDT, _ := sheetTitel.Cell(28, xlsx.ColLettersToIndex("D"))
	CellDTValue := strings.TrimSuffix(strings.TrimPrefix(CellDT.Value, "Направление подготовки "), "_x000D_\n")
	CellPP, _ := sheetTitel.Cell(29, xlsx.ColLettersToIndex("D"))
	guid_program, err := s.repo.CreateProgramm(CellPP.Value, CellDTValue)
	if err != nil {
		return err
	}
	if err := s.repo.CreateRelationship(guid_program, guid_faculty, "ProgrammFaculty"); err != nil {
		return err
	}
	CellFT, _ := sheetTitel.Cell(41, xlsx.ColLettersToIndex("C"))
	CellFT_value := strings.Split(CellFT.Value, " ")[2]
	CellQ, _ := sheetTitel.Cell(39, xlsx.ColLettersToIndex("C"))
	CellQ_value := strings.Split(CellQ.Value, " ")[1]
	CellG, _ := sheetTitel.Cell(39, xlsx.ColLettersToIndex("W"))
	CellNameF, _ := sheetTitel.Cell(37, xlsx.ColLettersToIndex("D"))

	sheetPlan, ok := wb.Sheet["План"]
	if !ok {
		return errors.New("Sheet not found")
	}

	max_rows := sheetPlan.MaxRow
	max_cols := sheetPlan.MaxCol

	col_start, err := get_start_col(max_cols, sheetPlan)
	if err != nil {
		return err
	}
	count_sermestr, err := counter_semestr(wb)
	if err != nil {
		return err
	}
	end_col, err := get_end_col(wb, max_cols)
	if err != nil {
		return err
	}
	d_z, err := get_d_z(max_cols, col_start, wb)
	if err != nil {
		return err
	}
	_row, err := get_start_row(max_rows, sheetPlan)
	if err != nil {
		return err
	}

	for _row := _row; _row < max_rows; _row++ {
		cell, _ := sheetPlan.Cell(_row, end_col-1)
		if cell.Value != "" {

			final_dict := make(map[string]string)
			fmt.Println(CellDT.Value)
			final_dict[key_direction_training] = CellDTValue
			final_dict[key_training_profile] = CellPP.Value
			final_dict[key_form_training] = CellFT_value
			final_dict[key_qualification] = CellQ_value
			final_dict[key_year_admission] = CellG.Value

			cellName, _ := sheetPlan.Cell(_row, 2)
			final_dict[key_name] = cellName.Value
			logrus.Info(final_dict[key_name])
			cellCode, _ := sheetPlan.Cell(_row, 1)
			final_dict[key_code] = cellCode.Value
			CellCompit, _ := sheetPlan.Cell(_row, end_col)
			str1, str2 := get_rezult_str(CellCompit.Value, dict_competencies)
			final_dict[key_purpose_mastering] = str1
			final_dict[key_results_mastering] = str2
			mas_places_dissiple := []string{"Дисциплина относится к обязательной части образовательной программы.", "Дисциплина относится к части образовательной программы, формируемой участниками образовательных отношений, является обязательной для изучения.", "Дисциплина относится к части образовательной программы, формируемой участниками образовательных отношений, предлагается обучающимся на выбор."}
			if strings.Contains(final_dict[key_code], "О") {
				final_dict[key_place_of_discipline] = mas_places_dissiple[0]
			} else {
				if strings.Contains(final_dict[key_code], "ДВ") || strings.Contains(final_dict[key_code], "ФТД") {
					final_dict[key_place_of_discipline] = mas_places_dissiple[2]
				} else {
					final_dict[key_place_of_discipline] = mas_places_dissiple[1]
				}
			}

			str_control, _ := get_form_control(_row, sheetPlan)
			final_dict[key_semester_development] = str_control
			cellZE, _ := sheetPlan.Cell(_row, xlsx.ColLettersToIndex("H"))
			final_dict[key_ZE] = cellZE.Value
			cellK, _ := sheetPlan.Cell(_row, xlsx.ColLettersToIndex("K"))
			final_dict[key_K] = cellK.Value
			final_dict[key_L] = get_sum(_row, max_cols, count_sermestr, d_z, wb, sheetPlan, "Лек")
			final_dict[key_C] = get_sum(_row, max_cols, count_sermestr, d_z, wb, sheetPlan, "Сем")
			final_dict[key_P] = get_sum(_row, max_cols, count_sermestr, d_z, wb, sheetPlan, "Пр")
			final_dict[key_LR] = get_sum(_row, max_cols, count_sermestr, d_z, wb, sheetPlan, "Лаб")

			if IsEngByLoop(final_dict[key_training_profile]) {
				final_dict[key_ImplementationLanguage] = "Английский"
			} else {
				final_dict[key_ImplementationLanguage] = "Русский"
			}
			final_dict[key_NameFacylte] = CellNameF.Value

			mas_key := []string{
				"DevelopmentMastering",
				"EntranceRequirements",
				"ContentDiscipline",
				"CurrentControl",
				"EvaluationProcedure",
				"MethodologicalSupport",
				"References",
				"ListInformationTechnologies",
				"MaterialSupport",
				"InformationDevelopers"}

			for _, value := range mas_key {
				final_dict[value] = ""
			}

			guid_node, err := s.repo.CreatePlan(final_dict)
			if err != nil {
				return err
			}
			//fmt.Println("CREATE " + guid)
			if err := s.repo.CreateRelationship(guid_node, guid_program, "PlanProgramm"); err != nil {
				return err
			}

			if err := s.repo.CreateRelationship(guid_node+"_html", guid_node, "HTMLTextPlan"); err != nil {
				return err
			}
		}
	}

	logrus.Info("Finish CreatePlans")
	return nil
}

func get_d_z(max_col int, start_col int, wb *xlsx.File) (int, error) {
	sheet, ok := wb.Sheet["План"]
	if !ok {
		return 0, errors.New("Sheet not found")
	}
	counter := 1
	_row := 0

	flag, err := get_flag_table_semestr(wb)
	if err != nil {
		return 0, err
	}
	if flag {
		_row = 2
	} else {
		_row = 1
	}

	cell, err := sheet.Cell(_row, start_col)
	if err != nil {
		return 0, err
	}
	for col := start_col + 1; col < max_col; col++ {
		cellTest, err := sheet.Cell(_row, col)
		if err != nil {
			return 0, err
		}
		if cellTest.Value != cell.Value {
			counter++
		} else {
			break
		}
	}
	return counter, nil

}

func get_start_col(max_col int, sheet *xlsx.Sheet) (int, error) {
	for _col := 0; _col <= max_col; _col++ {
		cell, err := sheet.Cell(0, _col)
		if err != nil {
			return 0, err
		}
		if cell.Value == "Курс 1" {
			return int(_col), nil
		}
	}
	return 0, errors.New("That went wrong")
}

func get_start_row(max_row int, sheet *xlsx.Sheet) (int, error) {
	for _row := 0; _row <= max_row; _row++ {
		cell, err := sheet.Cell(_row, 0)
		if err != nil {
			return 0, err
		}
		if cell.Value == "+" {
			return int(_row), nil
		}
	}
	return 0, errors.New("That went wrong")
}

func get_flag_table_semestr(wb *xlsx.File) (bool, error) {
	sheet, ok := wb.Sheet["План"]
	if !ok {
		return false, errors.New("Sheet not found")
	}
	cell, err := sheet.Cell(1, 0)
	if err != nil {
		return false, err
	}
	if cell.Value == "" {
		return true, nil
	} else {
		return false, nil
	}
}

func counter_semestr(wb *xlsx.File) (int, error) {
	sheet, ok := wb.Sheet["Титул"]
	if !ok {
		return 0, errors.New("Sheet not found")
	}
	cell, err := sheet.Cell(42, xlsx.ColLettersToIndex("C"))
	if err != nil {
		return 0, err
	}
	var term int
	if strings.Split(cell.Value, ":")[1][2] == '0' {
		term = 10
	} else {
		term = int(strings.Split(cell.Value, ":")[1][1])
	}
	flag, err := get_flag_table_semestr(wb)
	if err != nil {
		return 0, err
	}
	if flag {
		return term * 2, nil
	} else {
		return term, nil
	}
}

func get_end_col(wb *xlsx.File, max_col int) (int, error) {
	sheet, ok := wb.Sheet["План"]
	if !ok {
		return 0, errors.New("Sheet not found")
	}
	flag, err := get_flag_table_semestr(wb)
	if err != nil {
		return 0, err
	}
	if flag {
		for _col := 0; _col < max_col+1; _col++ {
			cell, err := sheet.Cell(2, _col)
			if err != nil {
				return 0, err
			}
			if cell.Value == "Компетенции" {
				return _col, nil
			}
		}
	} else {
		for _col := 0; _col < max_col+1; _col++ {
			cell, err := sheet.Cell(1, _col)
			if err != nil {
				return 0, err
			}
			if cell.Value == "Компетенции" {
				return _col, nil
			}
		}
	}
	return 0, errors.New("That went wrong")
}

func get_rezult_str(str_mas_cod string, dict_competencies map[string]dictCompetenciesStruct) (string, string) {
	mas_str := strings.Split(str_mas_cod, ";")
	for index, cod := range mas_str {
		if cod[0] == ' ' {
			mas_str[index] = mas_str[index][1:len(mas_str[index])]
		}
	}

	str1 := ""
	str2 := ""

	dict_competencies_final := make(map[string][]string)
	for key, value := range dict_competencies {
		for key2, value2 := range value.dictComp {
			dict_competencies_final[key2] = []string{value2, key + " " + value.title}
		}
	}

	for index, key := range mas_str {
		str2 += "\t" + key + " " + dict_competencies_final[key][0]
		if index != len(mas_str)-1 {
			str2 += "\n"
		}
		if !strings.Contains(str1, dict_competencies_final[key][1]) {
			str1 += "\t" + dict_competencies_final[key][1]
			if index != len(mas_str)-1 {
				str1 += "\n"
			}
		}
	}

	return str1, str2
}

func get_form_control(_row int, sheet *xlsx.Sheet) (string, error) {
	final_str := ""
	mas_form_control := []string{"Экзамен", "Зачет", "Зачет с оценкой"}
	dict := make(map[int]string)

	for index := 3; index < 6; index++ {
		cell, err := sheet.Cell(_row, index)
		if err != nil {
			return "", err
		}
		all_semester := cell.Value
		if all_semester != "" {
			for _, semester := range all_semester {
				int_sem, _ := strconv.Atoi(string(semester))
				dict[int_sem] = "Семестр " + string(semester) + ", " + mas_form_control[index-3]
			}
		}
	}

	keys := make([]int, 0, len(dict))
	for k := range dict {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	for _, k := range keys {
		final_str += dict[k] + "\n"
	}

	return final_str[:len(final_str)-1], nil
}

func get_sum(_row int, max_col int, counter_semestr int, d_z int, wb *xlsx.File, sheet *xlsx.Sheet, key string) string {
	sum := 0.0
	var _col int
	flag, _ := get_flag_table_semestr(wb)
	if flag {
		for col := 0; col <= max_col; col++ {
			cell, _ := sheet.Cell(2, col)
			if cell.Value == key {
				_col = col
				break
			}
		}
	} else {
		for col := 0; col <= max_col; col++ {
			cell, _ := sheet.Cell(1, col)
			if cell.Value == key {
				_col = col
				break
			}
		}
	}

	for i := 0; i < counter_semestr; i++ {
		cell, _ := sheet.Cell(_row, _col)
		value := cell.Value
		if value != "" {
			s, _ := strconv.ParseFloat(value, 64)
			sum += s
		}
		_col += d_z
	}

	return fmt.Sprint(sum)
}

func get_dict_competencies(wb *xlsx.File) (map[string]dictCompetenciesStruct, error) {
	logrus.Info("Start Proccesing")
	dict_competencies := make(map[string]dictCompetenciesStruct)
	if err := postProccesingWB(wb); err != nil {
		return nil, err
	}
	logrus.Info("Finish Proccesing")
	sheet, ok := wb.Sheet["Компетенции"]
	if !ok {
		return nil, errors.New("Sheet not found")
	}

	row := 1
	key := ""
	for {
		cell, err := sheet.Cell(row, 4)
		if err != nil {
			return nil, err
		}
		style := cell.GetStyle()
		if !style.ApplyBorder {
			break
		}

		if cell.Value != "" && cell.Value != "-" {
			rowTest, err := sheet.Row(row)
			if err != nil {
				return nil, err
			}
			var masTest []string
			rowTest.ForEachCell(func(c *xlsx.Cell) error {
				value := c.Value
				masTest = append(masTest, value)
				return nil
			})
			for _, value := range masTest {

				if value != "" {

					key = value
					cellTitle, err := sheet.Cell(row, 3)
					if err != nil {
						return nil, err
					}
					dict_competencies[key] = dictCompetenciesStruct{
						title:    cellTitle.Value,
						dictComp: make(map[string]string),
					}
					break
				}
			}

		} else if cell.Value == "-" {
			rowTest, err := sheet.Row(row)
			if err != nil {
				return nil, err
			}
			var masTest []string
			rowTest.ForEachCell(func(c *xlsx.Cell) error {
				value := c.Value
				masTest = append(masTest, value)
				return nil
			})
			for _, value := range masTest {
				if value != "" {
					cellDrop, err := sheet.Cell(row, 3)
					if err != nil {
						return nil, err
					}
					dict_competencies[key].dictComp[value] = cellDrop.Value
					break
				}
			}
		}
		row++
	}
	return dict_competencies, nil
}

func postProccesingWB(wb *xlsx.File) error {
	sheet, ok := wb.Sheet["Компетенции"]
	if !ok {
		return errors.New("Sheet not found")
	}
	row := 1
	for {
		cell, err := sheet.Cell(row, 4)
		if err != nil {
			return err
		}
		style := cell.GetStyle()
		if !style.ApplyBorder {
			break
		}

		rowTest, err := sheet.Row(row)
		if err != nil {
			return err
		}
		var masTest []string
		rowTest.ForEachCell(func(c *xlsx.Cell) error {
			value := c.Value
			masTest = append(masTest, value)
			return nil
		})

		for _, value := range masTest {
			if value != "" {
				switch strings.Contains(value, "-") {
				case true:
					reg, err := regexp.MatchString(`\S\d`, value)
					if err != nil {
						return err
					}
					if reg {
						cell.SetValue(strings.Split(value, "-")[0])
					}
				case false:
					reg, err := regexp.MatchString(`И.*\d+`, value)
					if err != nil {
						return err
					}
					if reg {
						cell.SetValue("-")
					}
				}
				break
			}
		}

		row++
	}
	return nil
}
