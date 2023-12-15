package myDes

import (
	"fmt"
	"strconv"
)

// MyDES представляет собой алгоритм шифрования/дешифрования DES
type MyDES struct {
	childKeys []string // Массив для хранения подключей
	iv        string   // Вектор инициализации
}

// NewMyDES инициализирует новый экземпляр MyDES с заданным вектором инициализации
func NewMyDES(iv string) *MyDES {
	return &MyDES{
		iv: iv,
	}
}

// bitEncode преобразует строку в бинарное представление (01)
func (d *MyDES) bitEncode(s string) string {
	binStr := ""
	for _, c := range []byte(s) {
		// Преобразование каждого символа в 8-битное бинарное представление и добавление к binStr
		binStr += fmt.Sprintf("%08b", c)
	}
	return binStr
}

// bitDecode преобразует список бинарных строк обратно в исходную строку
func (d *MyDES) bitDecode(s []string) string {
	decoded := ""
	for _, binStr := range s {
		// Преобразование каждой бинарной строки в целое число и далее в символ
		val, _ := strconv.ParseInt(binStr, 2, 64)
		decoded += string(val)
	}
	return decoded
}

// negate инвертирует бинарную строку
func (d *MyDES) negate(s string) string {
	result := ""
	for _, i := range s {
		// Инверсия каждого бита в бинарной строке
		if i == '1' {
			result += "0"
		} else {
			result += "1"
		}
	}
	return result
}

// replaceBlock заменяет отдельные биты блока на основе таблицы замены
func (d *MyDES) replaceBlock(block string, replaceTable []int) string {
	result := ""
	for _, i := range replaceTable {
		// Замена битов в блоке согласно указанным позициям в таблице замены
		result += string(block[i-1])
	}
	return result
}

// processingEncodeInput преобразует входную строку в бинарную форму и разбивает ее на блоки по 64 бита
func (d *MyDES) processingEncodeInput(input string) []string {
	result := make([]string, 0)
	bitString := d.bitEncode(input)

	// Если длина бинарной строки не кратна 64, добавляем нули для выравнивания
	if len(bitString)%64 != 0 {
		for i := 0; i < 64-len(bitString)%64; i++ {
			bitString += "0"
		}
	}

	// Разбиваем бинарную строку на блоки по 64 бита
	for i := 0; i < len(bitString)/64; i++ {
		result = append(result, bitString[i*64:i*64+64])
	}
	return result
}

// processingDecodeInput преобразует входную строку в шестнадцатеричной форме в бинарную и разбивает ее на блоки по 64 бита
func (d *MyDES) processingDecodeInput(input []byte) []string {
	result := make([]string, 0)
	inputList := make([]string, 0)

	// Конвертация каждого байта в шестнадцатеричное представление и добавление в inputList
	for _, b := range input {
		inputList = append(inputList, fmt.Sprintf("%02x", b))
	}

	intList := make([]int64, 0)

	// Конвертация шестнадцатеричных строк в целые числа
	for _, i := range inputList {
		val, _ := strconv.ParseInt(i, 16, 64)
		intList = append(intList, val)
	}

	// Конвертация целых чисел в бинарные строки, добавление в результат и выравнивание до 64 бит
	for _, i := range intList {
		binData := strconv.FormatInt(i, 2)
		for len(binData) < 64 {
			binData = "0" + binData
		}
		result = append(result, binData)
	}
	return result
}

// keyConversion преобразует исходный 64-битный ключ в 56-битный ключ и выполняет замену
func (d *MyDES) keyConversion(key string) string {
	// Преобразование ключа в бинарную строку
	key = d.bitEncode(key)

	// Добавление нулей до достижения 64 бит
	for len(key) < 64 {
		key += "0"
	}
	firstKey := key[:64]

	// Таблица для замены битов в ключе
	keyReplaceTable := []int{
		57, 49, 41, 33, 25, 17, 9, 1, 58, 50, 42, 34, 26, 18,
		10, 2, 59, 51, 43, 35, 27, 19, 11, 3, 60, 52, 44, 36,
		63, 55, 47, 39, 31, 23, 15, 7, 62, 54, 46, 38, 30, 22,
		14, 6, 61, 53, 45, 37, 29, 21, 13, 5, 28, 20, 12, 4,
	}

	// Замена битов в ключе согласно таблице
	return d.replaceBlock(firstKey, keyReplaceTable)
}

// spinKey выполняет вращение для генерации подключей
func (d *MyDES) spinKey(key string) []string {
	// Получение 56-битного ключа после замены
	kc := d.keyConversion(key)
	first, second := kc[:28], kc[28:]

	// Таблица для вращения ключа
	spinTable := []int{1, 2, 4, 6, 8, 10, 12, 14, 15, 17, 19, 21, 23, 25, 27, 28}
	subKeys := make([]string, 16)

	// Выполнение вращения и создание 16 подключей
	for i := 0; i < 16; i++ {
		firstAfterSpin := first[spinTable[i]:] + first[:spinTable[i]]
		secondAfterSpin := second[spinTable[i]:] + second[:spinTable[i]]
		subKeys[i] = firstAfterSpin + secondAfterSpin
	}
	return subKeys
}

// keySelectionReplacement получает подключ в 48 бит путем выборочной перестановки
func (d *MyDES) keySelectionReplacement(key string) {
	// Обнуление массива подключей
	d.childKeys = nil

	// Таблица для выборочной перестановки битов в подключе
	keySelectTable := []int{
		14, 17, 11, 24, 1, 5, 3, 28, 15, 6, 21, 10,
		23, 19, 12, 4, 26, 8, 16, 7, 27, 20, 13, 2,
		41, 52, 31, 37, 47, 55, 30, 40, 51, 45, 33, 48,
		44, 49, 39, 56, 34, 53, 46, 42, 50, 36, 29, 32,
	}

	// Генерация подключей
	subKeys := d.spinKey(key)

	// Добавление подключей после выборочной перестановки
	for _, childKey56 := range subKeys {
		d.childKeys = append(d.childKeys, d.replaceBlock(childKey56, keySelectTable))
	}
}

// initReplaceBlock выполняет начальную блочную перестановку
func (d *MyDES) initReplaceBlock(block string) string {
	// Таблица для начальной блочной перестановки
	replaceTable := []int{
		58, 50, 42, 34, 26, 18, 10, 2,
		60, 52, 44, 36, 28, 20, 12, 4,
		62, 54, 46, 38, 30, 22, 14, 6,
		64, 56, 48, 40, 32, 24, 16, 8,
		57, 49, 41, 33, 25, 17, 9, 1,
		59, 51, 43, 35, 27, 19, 11, 3,
		61, 53, 45, 37, 29, 21, 13, 5,
		63, 55, 47, 39, 31, 23, 15, 7,
	}

	// Выполнение блочной перестановки
	return d.replaceBlock(block, replaceTable)
}

// endReplaceBlock выполняет конечную блочную перестановку
func (d *MyDES) endReplaceBlock(block string) string {
	// Таблица для конечной блочной перестановки
	replaceTable := []int{
		40, 8, 48, 16, 56, 24, 64, 32,
		39, 7, 47, 15, 55, 23, 63, 31,
		38, 6, 46, 14, 54, 22, 62, 30,
		37, 5, 45, 13, 53, 21, 61, 29,
		36, 4, 44, 12, 52, 20, 60, 28,
		35, 3, 43, 11, 51, 19, 59, 27,
		34, 2, 42, 10, 50, 18, 58, 26,
		33, 1, 41, 9, 49, 17, 57, 25,
	}

	// Выполнение блочной перестановки
	return d.replaceBlock(block, replaceTable)
}

// blockExtend расширяет блок с использованием расширения
func (d *MyDES) blockExtend(block string) string {
	extendedBlock := ""   // Инициализируем пустую строку, в которую будем добавлять расширенные биты блока
	extendTable := []int{ // Таблица расширения, определяющая порядок выбора битов из блока
		32, 1, 2, 3, 4, 5,
		4, 5, 6, 7, 8, 9,
		8, 9, 10, 11, 12, 13,
		12, 13, 14, 15, 16, 17,
		16, 17, 18, 19, 20, 21,
		20, 21, 22, 23, 24, 25,
		24, 25, 26, 27, 28, 29,
		28, 29, 30, 31, 32, 1,
	}

	// Проходим по каждому индексу в таблице и добавляем соответствующий бит из блока в расширенный блок
	for _, i := range extendTable {
		extendedBlock += string(block[i-1])
	}

	return extendedBlock // Возвращаем расширенный блок
}

// notOr выполняет операцию XOR двух бинарных строк (01)
func (d *MyDES) notOr(a, b string) string {
	size := len(a) // Получаем длину одной из строк, предполагая, что длины обеих строк равны
	result := ""   // Инициализируем строку для хранения результата XOR

	// Проходим по каждому биту в строках и выполняем XOR, добавляя результат в строку результата
	for i := 0; i < size; i++ {
		if a[i] == b[i] {
			result += "0"
		} else {
			result += "1"
		}
	}

	return result // Возвращаем результат операции XOR
}

// sBoxReplace выполняет подстановку S-Box, преобразуя входные 48 бит в выходные 32 бита
func (d *MyDES) sBoxReplace(block48 string) string {
	// Таблица S-Box, представленная как двумерный массив
	sBoxTable := [8][4][16]int{
		{
			{14, 4, 13, 1, 2, 15, 11, 8, 3, 10, 6, 12, 5, 9, 0, 7},
			{0, 15, 7, 4, 14, 2, 13, 1, 10, 6, 12, 11, 9, 5, 3, 8},
			{4, 1, 14, 8, 13, 6, 2, 11, 15, 12, 9, 7, 3, 10, 5, 0},
			{15, 12, 8, 2, 4, 9, 1, 7, 5, 11, 3, 14, 10, 0, 6, 13},
		},
		{
			{15, 1, 8, 14, 6, 11, 3, 4, 9, 7, 2, 13, 12, 0, 5, 10},
			{3, 13, 4, 7, 15, 2, 8, 14, 12, 0, 1, 10, 6, 9, 11, 5},
			{0, 14, 7, 11, 10, 4, 13, 1, 5, 8, 12, 6, 9, 3, 2, 15},
			{13, 8, 10, 1, 3, 15, 4, 2, 11, 6, 7, 12, 0, 5, 14, 9},
		},
		{
			{10, 0, 9, 14, 6, 3, 15, 5, 1, 13, 12, 7, 11, 4, 2, 8},
			{13, 7, 0, 9, 3, 4, 6, 10, 2, 8, 5, 14, 12, 11, 15, 1},
			{13, 6, 4, 9, 8, 15, 3, 0, 11, 1, 2, 12, 5, 10, 14, 7},
			{1, 10, 13, 0, 6, 9, 8, 7, 4, 15, 14, 3, 11, 5, 2, 12},
		},
		{
			{7, 13, 14, 3, 0, 6, 9, 10, 1, 2, 8, 5, 11, 12, 4, 15},
			{13, 8, 11, 5, 6, 15, 0, 3, 4, 7, 2, 12, 1, 10, 14, 9},
			{10, 6, 9, 0, 12, 11, 7, 13, 15, 1, 3, 14, 5, 2, 8, 4},
			{3, 15, 0, 6, 10, 1, 13, 8, 9, 4, 5, 11, 12, 7, 2, 14},
		},
		{
			{2, 12, 4, 1, 7, 10, 11, 6, 8, 5, 3, 15, 13, 0, 14, 9},
			{14, 11, 2, 12, 4, 7, 13, 1, 5, 0, 15, 10, 3, 9, 8, 6},
			{4, 2, 1, 11, 10, 13, 7, 8, 15, 9, 12, 5, 6, 3, 0, 14},
			{11, 8, 12, 7, 1, 14, 2, 13, 6, 15, 0, 9, 10, 4, 5, 3},
		},
		{
			{12, 1, 10, 15, 9, 2, 6, 8, 0, 13, 3, 4, 14, 7, 5, 11},
			{10, 15, 4, 2, 7, 12, 9, 5, 6, 1, 13, 14, 0, 11, 3, 8},
			{9, 14, 15, 5, 2, 8, 12, 3, 7, 0, 4, 10, 1, 13, 11, 6},
			{4, 3, 2, 12, 9, 5, 15, 10, 11, 14, 1, 7, 6, 0, 8, 13},
		},
		{
			{4, 11, 2, 14, 15, 0, 8, 13, 3, 12, 9, 7, 5, 10, 6, 1},
			{13, 0, 11, 7, 4, 9, 1, 10, 14, 3, 5, 12, 2, 15, 8, 6},
			{1, 4, 11, 13, 12, 3, 7, 14, 10, 15, 6, 8, 0, 5, 9, 2},
			{6, 11, 13, 8, 1, 4, 10, 7, 9, 5, 0, 15, 14, 2, 3, 12},
		},
		{
			{13, 2, 8, 4, 6, 15, 11, 1, 10, 9, 3, 14, 5, 0, 12, 7},
			{1, 15, 13, 8, 10, 3, 7, 4, 12, 5, 6, 11, 0, 14, 9, 2},
			{7, 11, 4, 1, 9, 12, 14, 2, 0, 6, 10, 13, 15, 3, 5, 8},
			{2, 1, 14, 7, 4, 10, 8, 13, 15, 12, 9, 0, 3, 5, 6, 11},
		},
	}

	result := "" // Инициализируем строку для хранения результата замены S-Box
	for i := 0; i < 8; i++ {
		// Получаем биты строки и столбца для текущего блока 6 бит
		rowBit := string([]byte{block48[i*6], block48[i*6+5]})
		lineBit := block48[i*6+1 : i*6+5]

		// Преобразуем строку и столбец в целочисленные значения
		row, _ := strconv.ParseInt(rowBit, 2, 64)
		line, _ := strconv.ParseInt(lineBit, 2, 64)

		// Получаем значение из таблицы S-Box и преобразуем его в бинарную строку
		data := sBoxTable[i][row][line]
		noFull := fmt.Sprintf("%04b", data)
		result += noFull
	}

	return result // Возвращаем результат замены S-Box
}

// sBoxCompression выполняет компрессию блока S-Box для 48-битного блока согласно таблице компрессии S-Box
func (d *MyDES) sBoxCompression(num int, block48 string) string {
	// Выполняем операцию NOT OR между 48-битным блоком и подключом
	resultNotOr := d.notOr(block48, d.childKeys[num])

	// Выполняем подстановку S-Box, используя полученный результат
	return d.sBoxReplace(resultNotOr)
}

// pBoxReplacement заменяет 32-битный блок с использованием таблицы замены P-Box
func (d *MyDES) pBoxReplacement(block32 string) string {
	// Таблица замены P-Box для DES
	pBoxReplaceTable := []int{
		16, 7, 20, 21, 29, 12, 28, 17, 1, 15, 23, 26, 5, 18, 31, 10,
		2, 8, 24, 14, 32, 27, 3, 9, 19, 13, 30, 6, 22, 11, 4, 25,
	}

	// Выполняем замену P-Box
	return d.replaceBlock(block32, pBoxReplaceTable)
}

// fFunction представляет собой функцию F, часть сети Фейстеля
func (d *MyDES) fFunction(right string, isDecode bool, num int) string {
	// Убедимся, что правая половина блока расширена до 48 бит
	right = d.blockExtend(right)

	var sbcResult string
	if isDecode {
		// Для расшифровки выполняем компрессию блока S-Box с конкретным подключом
		sbcResult = d.sBoxCompression(15-num, right)
	} else {
		// Для шифрования выполняем компрессию блока S-Box с конкретным подключом
		sbcResult = d.sBoxCompression(num, right)
	}

	// Выполняем замену P-Box для результата S-Box
	return d.pBoxReplacement(sbcResult)
}

// iteration выполняет одну итерацию сети Фейстеля
func (d *MyDES) iteration(block string, key string, isDecode bool) string {
	// Выбираем подключи на основе ключа
	d.keySelectionReplacement(key)

	// Итерируем 16 раз (количество раундов в DES)
	for i := 0; i < 16; i++ {
		// Разбиваем блок на левую и правую половины
		left, right := block[0:32], block[32:64]

		// Сохраняем левую половину для следующей итерации
		nextLeft := right

		// Выполняем функцию F на правой половине
		fResult := d.fFunction(right, isDecode, i)

		// Выполняем операцию NOT OR между левой половиной и результатом F
		right = d.notOr(left, fResult)

		// Обновляем блок для следующей итерации
		block = nextLeft + right
	}

	// Сцепляем правую и левую половины блока и возвращаем результат
	return block[32:] + block[:32]
}

// encode выполняет шифрование DES в режиме CBC
func (d *MyDES) Encode(input string, key string) string {
	// Результирующая строка для закодированных блоков
	result := ""

	// Обрабатываем входную строку для получения блоков
	blocks := d.processingEncodeInput(input)

	// Используем IV как предыдущий блок для первой итерации
	previousBlock := d.bitEncode(d.iv)

	// Итерируем по блокам
	for _, block := range blocks {
		// Выполняем операцию NOT OR с предыдущим блоком
		block = d.notOr(block, previousBlock)

		// Заменяем блок перед итерацией
		irbResult := d.initReplaceBlock(block)

		// Выполняем одну итерацию сети Фейстеля
		blockResult := d.iteration(irbResult, key, false)

		// Заменяем блок после итерации
		blockResult = d.endReplaceBlock(blockResult)

		// Преобразуем результат в шестнадцатеричную форму и добавляем к общему результату
		result += fmt.Sprintf("%x", blockResult)

		// Обновляем предыдущий блок для следующей итерации
		previousBlock = blockResult
	}

	return result
}

// decode выполняет расшифровку DES в режиме CBC
func (d *MyDES) Decode(cipherText []byte, key string) string {
	// Результирующий массив для расшифрованных блоков
	var result []string

	// Обрабатываем входные данные для получения блоков
	blocks := d.processingDecodeInput(cipherText)

	// Используем IV как предыдущий блок для первой итерации
	previousBlock := d.bitEncode(d.iv)

	// Итерируем по блокам
	for _, block := range blocks {
		// Заменяем блок перед итерацией
		irbResult := d.initReplaceBlock(block)

		// Выполняем одну итерацию сети Фейстеля для расшифровки
		blockResult := d.iteration(irbResult, key, true)

		// Заменяем блок после итерации
		blockResult = d.endReplaceBlock(blockResult)

		// Выполняем операцию NOT OR с предыдущим блоком
		blockResult = d.notOr(blockResult, previousBlock)

		// Разбиваем результат на 8-битные части и добавляем к общему результату
		for i := 0; i < len(blockResult); i += 8 {
			result = append(result, blockResult[i:i+8])
		}

		// Обновляем предыдущий блок для следующей итерации
		previousBlock = block
	}

	// Преобразуем конечный результат в строку
	return d.bitDecode(result)
}
