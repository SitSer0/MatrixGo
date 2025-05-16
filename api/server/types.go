package server

// MatrixRequest представляет запрос с матрицей
type MatrixRequest struct {
	Type string     `json:"type"`           // "float64", "complex", "rational", "gf"
	Rows int        `json:"rows"`           // Количество строк
	Cols int        `json:"cols"`           // Количество столбцов
	Data [][]string `json:"data"`           // Значения в строковом формате
	ModP int64      `json:"modP,omitempty"` // Для конечного поля GF(p)
}

// SystemRequest представляет запрос для решения системы уравнений
type SystemRequest struct {
	Matrix MatrixRequest `json:"matrix"` // Матрица системы
	Vector []string      `json:"vector"` // Вектор правой части
}

// MatrixResponse представляет ответ сервера
type MatrixResponse struct {
	Result [][]string `json:"result,omitempty"` // Для матричных результатов
	Value  string     `json:"value,omitempty"`  // Для скалярных результатов
	Error  string     `json:"error,omitempty"`  // Сообщение об ошибке
}
