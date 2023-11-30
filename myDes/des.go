package myDes

import (
	"fmt"
	"strconv"
)

// MyDES представляет собой алгоритм шифрования/дешифрования DES
type MyDES struct {
	childKeys []string
	iv        string
}

// NewMyDES инициализирует новый экземпляр MyDES
func NewMyDES(iv string) *MyDES {
	return &MyDES{
		iv: iv,
	}
}

// bitEncode преобразует строку в бинарное представление (01)
func (d *MyDES) bitEncode(s string) string {
	binStr := ""
	for _, c := range []byte(s) {
		binStr += fmt.Sprintf("%08b", c)
	}
	return binStr
}

// bitDecode преобразует список бинарных строк обратно в исходную строку
func (d *MyDES) bitDecode(s []string) string {
	decoded := ""
	for _, binStr := range s {
		val, _ := strconv.ParseInt(binStr, 2, 64)
		decoded += string(val)
	}
	return decoded
}

// negate инвертирует бинарную строку
func (d *MyDES) negate(s string) string {
	result := ""
	for _, i := range s {
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
		result += string(block[i-1])
	}
	return result
}

// processingEncodeInput преобразует входную строку в бинарную форму и разбивает ее на блоки по 64 бита
func (d *MyDES) processingEncodeInput(input string) []string {
	result := make([]string, 0)
	bitString := d.bitEncode(input)
	// Если длина не кратна 64, добавляем нули
	if len(bitString)%64 != 0 {
		for i := 0; i < 64-len(bitString)%64; i++ {
			bitString += "0"
		}
	}
	for i := 0; i < len(bitString)/64; i++ {
		result = append(result, bitString[i*64:i*64+64])
	}
	return result
}

// processingDecodeInput преобразует входную строку в шестнадцатеричной форме в бинарную и разбивает ее на блоки по 64 бита
func (d *MyDES) processingDecodeInput(input []byte) []string {
	result := make([]string, 0)
	inputList := make([]string, 0)
	for _, b := range input {
		inputList = append(inputList, fmt.Sprintf("%02x", b))
	}
	intList := make([]int64, 0)
	for _, i := range inputList {
		val, _ := strconv.ParseInt(i, 16, 64)
		intList = append(intList, val)
	}
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
	key = d.bitEncode(key)
	for len(key) < 64 {
		key += "0"
	}
	firstKey := key[:64]
	keyReplaceTable := []int{
		57, 49, 41, 33, 25, 17, 9, 1, 58, 50, 42, 34, 26, 18,
		10, 2, 59, 51, 43, 35, 27, 19, 11, 3, 60, 52, 44, 36,
		63, 55, 47, 39, 31, 23, 15, 7, 62, 54, 46, 38, 30, 22,
		14, 6, 61, 53, 45, 37, 29, 21, 13, 5, 28, 20, 12, 4,
	}
	return d.replaceBlock(firstKey, keyReplaceTable)
}

// spinKey выполняет вращение для генерации подключей
func (d *MyDES) spinKey(key string) []string {
	kc := d.keyConversion(key)
	first, second := kc[:28], kc[28:]
	spinTable := []int{1, 2, 4, 6, 8, 10, 12, 14, 15, 17, 19, 21, 23, 25, 27, 28}
	subKeys := make([]string, 16)
	for i := 0; i < 16; i++ {
		firstAfterSpin := first[spinTable[i]:] + first[:spinTable[i]]
		secondAfterSpin := second[spinTable[i]:] + second[:spinTable[i]]
		subKeys[i] = firstAfterSpin + secondAfterSpin
	}
	return subKeys
}

// keySelectionReplacement получает подключ в 48 бит путем выборочной перестановки
func (d *MyDES) keySelectionReplacement(key string) {
	d.childKeys = nil
	keySelectTable := []int{
		14, 17, 11, 24, 1, 5, 3, 28, 15, 6, 21, 10,
		23, 19, 12, 4, 26, 8, 16, 7, 27, 20, 13, 2,
		41, 52, 31, 37, 47, 55, 30, 40, 51, 45, 33, 48,
		44, 49, 39, 56, 34, 53, 46, 42, 50, 36, 29, 32,
	}
	subKeys := d.spinKey(key)
	for _, childKey56 := range subKeys {
		d.childKeys = append(d.childKeys, d.replaceBlock(childKey56, keySelectTable))
	}
}

// initReplaceBlock выполняет начальную блочную перестановку
func (d *MyDES) initReplaceBlock(block string) string {
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
	return d.replaceBlock(block, replaceTable)
}

// endReplaceBlock выполняет конечную блочную перестановку
func (d *MyDES) endReplaceBlock(block string) string {
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
	return d.replaceBlock(block, replaceTable)
}

// blockExtend расширяет блок с использованием расширения
func (d *MyDES) blockExtend(block string) string {
	extendedBlock := ""
	extendTable := []int{
		32, 1, 2, 3, 4, 5,
		4, 5, 6, 7, 8, 9,
		8, 9, 10, 11, 12, 13,
		12, 13, 14, 15, 16, 17,
		16, 17, 18, 19, 20, 21,
		20, 21, 22, 23, 24, 25,
		24, 25, 26, 27, 28, 29,
		28, 29, 30, 31, 32, 1,
	}
	for _, i := range extendTable {
		extendedBlock += string(block[i-1])
	}
	return extendedBlock
}

// notOr выполняет операцию XOR двух бинарных строк (01)
func (d *MyDES) notOr(a, b string) string {
	size := len(a)
	result := ""
	for i := 0; i < size; i++ {
		if a[i] == b[i] {
			result += "0"
		} else {
			result += "1"
		}
	}
	return result
}

// sBoxReplace выполняет подстановку S-Box, преобразуя входные 48 бит в выходные 32 бита
func (d *MyDES) sBoxReplace(block48 string, num int) string {
	// Таблицы замены S-Box для DES
	sBoxTable := [][][]int{
		// 8 S-Box'ов, каждый со 4 строками и 16 столбцами
		{
			{14, 4, 13, 1, 2, 15, 11, 8, 3, 10, 6, 12, 5, 9, 0, 7},
			// ... (другие строки для первого S-Box'а)
		},
		// ... (другие S-Box'ы)
	}

	// Результирующий 32-битный блок после подстановки S-Box
	result := ""

	// Итерируем по 8 S-Box'ам
	for i := 0; i < 8; i++ {
		// Извлекаем соответствующие биты для выбора строки и столбца
		rowBit := string([]byte{block48[i*6], block48[i*6+5]})
		lineBit := block48[i*6+1 : i*6+5]

		// Преобразуем двоичные биты строки и столбца в десятичные
		row, _ := strconv.ParseInt(rowBit, 2, 64)
		line, _ := strconv.ParseInt(lineBit, 2, 64)

		// Получаем значение из таблицы S-Box
		data := sBoxTable[i][row][line]

		// Преобразуем результат в двоичную форму и убеждаемся, что она состоит из 4 бит
		noFull := strconv.FormatInt(int64(data), 2)
		for len(noFull) < 4 {
			noFull = "0" + noFull
		}

		// Добавляем 4-битный результат к общему результату
		result += noFull
	}

	return result
}

// sBoxCompression выполняет компрессию блока S-Box для 48-битного блока согласно таблице компрессии S-Box
func (d *MyDES) sBoxCompression(num int, block48 string) string {
	// Выполняем операцию NOT OR между 48-битным блоком и подключом
	resultNotOr := d.notOr(block48, d.childKeys[num])

	// Выполняем подстановку S-Box, используя полученный результат
	return d.sBoxReplace(resultNotOr, num)
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
