package main

type LangID int

const (
	EN LangID = iota
	TR
	ZH
	JA
	KO
	DE
	PT
	ES
	HI
	CA
	RU
	AZ
	KK
	UZ
)

type LangPack struct {
	WindowTitle       string
	NewDownload       string
	URL               string
	Connections       string
	Download          string
	CancelSelected    string
	ClearCompleted    string
	Filename          string
	Size              string
	Progress          string
	Speed             string
	Time              string
	ETA               string
	Status            string
	Conns             string
	Total             string
	Active            string
	Language          string
	Queued            string
	Downloading       string
	Completed         string
	Error             string
	Cancelled         string
	EnterURL          string
	InvalidURL        string
	ErrorTitle        string
	UnknownFileSize   string
	File              string
	Done              string
	StatusBarTotal    string
	StatusBarActive   string
	StatusBarSpeed    string
	ParallelConns     string
	SingleConnection  string
}

var currentLang = EN

var langs = map[LangID]LangPack{
	EN: {
		WindowTitle: "Fast Download Manager", NewDownload: "New Download",
		URL: "URL:", Connections: "Conns:", Download: "Download",
		CancelSelected: "Cancel Selected", ClearCompleted: "Clear Completed",
		Filename: "Filename", Size: "Size", Progress: "Progress", Speed: "Speed",
		Time: "Time", ETA: "ETA", Status: "Status", Conns: "Conns",
		Total: "Total:", Active: "Active:", Language: "Language",
		Queued: "Queued", Downloading: "Downloading", Completed: "Completed",
		Error: "Error", Cancelled: "Cancelled",
		EnterURL: "Please enter a URL", InvalidURL: "Invalid URL",
		ErrorTitle: "Error", UnknownFileSize: "Unknown file size",
		File: "File", Done: "Done",
		StatusBarTotal: "Total:", StatusBarActive: "Active:", StatusBarSpeed: "Speed:",
		ParallelConns: "parallel connections", SingleConnection: "single connection",
	},
	TR: {
		WindowTitle: "Hızlı İndirme Yöneticisi", NewDownload: "Yeni İndirme",
		URL: "Adres:", Connections: "Bağlnt:", Download: "İndir",
		CancelSelected: "Seçileni İptal Et", ClearCompleted: "Bitmişleri Temizle",
		Filename: "Dosya Adı", Size: "Boyut", Progress: "İlerleme", Speed: "Hız",
		Time: "Süre", ETA: "Kalan", Status: "Durum", Conns: "Bağlt",
		Total: "Toplam:", Active: "Aktif:", Language: "Dil",
		Queued: "Sırada", Downloading: "İndiriliyor", Completed: "Tamamlandı",
		Error: "Hata", Cancelled: "İptal Edildi",
		EnterURL: "Lütfen bir adres girin", InvalidURL: "Geçersiz adres",
		ErrorTitle: "Hata", UnknownFileSize: "Dosya boyutu bilinmiyor",
		File: "Dosya", Done: "Tamam",
		StatusBarTotal: "Toplam:", StatusBarActive: "Aktif:", StatusBarSpeed: "Hız:",
		ParallelConns: "paralel bağlantı", SingleConnection: "tek bağlantı",
	},
	ZH: {
		WindowTitle: "高速下载管理器", NewDownload: "新建下载",
		URL: "网址:", Connections: "连接:", Download: "下载",
		CancelSelected: "取消选中", ClearCompleted: "清除已完成",
		Filename: "文件名", Size: "大小", Progress: "进度", Speed: "速度",
		Time: "时间", ETA: "预计", Status: "状态", Conns: "连",
		Total: "总计:", Active: "活跃:", Language: "语言",
		Queued: "排队中", Downloading: "下载中", Completed: "已完成",
		Error: "错误", Cancelled: "已取消",
		EnterURL: "请输入网址", InvalidURL: "无效网址",
		ErrorTitle: "错误", UnknownFileSize: "未知文件大小",
		File: "文件", Done: "完成",
		StatusBarTotal: "总计:", StatusBarActive: "活跃:", StatusBarSpeed: "速度:",
		ParallelConns: "并行连接", SingleConnection: "单连接",
	},
	JA: {
		WindowTitle: "高速ダウンロード管理", NewDownload: "新規ダウンロード",
		URL: "URL:", Connections: "接続:", Download: "ダウンロード",
		CancelSelected: "選択をキャンセル", ClearCompleted: "完了を消去",
		Filename: "ファイル名", Size: "サイズ", Progress: "進捗", Speed: "速度",
		Time: "時間", ETA: "残り", Status: "状態", Conns: "接続",
		Total: "合計:", Active: "アクティブ:", Language: "言語",
		Queued: "待機中", Downloading: "ダウンロード中", Completed: "完了",
		Error: "エラー", Cancelled: "キャンセル済",
		EnterURL: "URLを入力してください", InvalidURL: "無効なURL",
		ErrorTitle: "エラー", UnknownFileSize: "ファイルサイズ不明",
		File: "ファイル", Done: "完了",
		StatusBarTotal: "合計:", StatusBarActive: "アクティブ:", StatusBarSpeed: "速度:",
		ParallelConns: "並列接続", SingleConnection: "単一接続",
	},
	KO: {
		WindowTitle: "고속 다운로드 관리자", NewDownload: "새 다운로드",
		URL: "URL:", Connections: "연결:", Download: "다운로드",
		CancelSelected: "선택 취소", ClearCompleted: "완료 지우기",
		Filename: "파일명", Size: "크기", Progress: "진행", Speed: "속도",
		Time: "시간", ETA: "예상", Status: "상태", Conns: "연결",
		Total: "합계:", Active: "활성:", Language: "언어",
		Queued: "대기중", Downloading: "다운로드중", Completed: "완료",
		Error: "오류", Cancelled: "취소됨",
		EnterURL: "URL을 입력하세요", InvalidURL: "잘못된 URL",
		ErrorTitle: "오류", UnknownFileSize: "알 수 없는 파일 크기",
		File: "파일", Done: "완료",
		StatusBarTotal: "합계:", StatusBarActive: "활성:", StatusBarSpeed: "속도:",
		ParallelConns: "병렬 연결", SingleConnection: "단일 연결",
	},
	DE: {
		WindowTitle: "Schnell Download-Manager", NewDownload: "Neuer Download",
		URL: "URL:", Connections: "Verb.:", Download: "Download",
		CancelSelected: "Auswahl abbrechen", ClearCompleted: "Erledigte löschen",
		Filename: "Dateiname", Size: "Größe", Progress: "Fortschritt", Speed: "Geschw.",
		Time: "Zeit", ETA: "Rest", Status: "Status", Conns: "Verb.",
		Total: "Gesamt:", Active: "Aktiv:", Language: "Sprache",
		Queued: "Wartend", Downloading: "Lädt", Completed: "Fertig",
		Error: "Fehler", Cancelled: "Abgebrochen",
		EnterURL: "Bitte URL eingeben", InvalidURL: "Ungültige URL",
		ErrorTitle: "Fehler", UnknownFileSize: "Unbekannte Dateigröße",
		File: "Datei", Done: "Erledigt",
		StatusBarTotal: "Gesamt:", StatusBarActive: "Aktiv:", StatusBarSpeed: "Geschw.:",
		ParallelConns: "parallele Verbindungen", SingleConnection: "Einzelverbindung",
	},
	PT: {
		WindowTitle: "Gerenciador de Downloads Rápido", NewDownload: "Novo Download",
		URL: "URL:", Connections: "Conexões:", Download: "Baixar",
		CancelSelected: "Cancelar Selecionado", ClearCompleted: "Limpar Concluídos",
		Filename: "Arquivo", Size: "Tamanho", Progress: "Progresso", Speed: "Veloc.",
		Time: "Tempo", ETA: "Restante", Status: "Estado", Conns: "Conex.",
		Total: "Total:", Active: "Ativo:", Language: "Idioma",
		Queued: "Na fila", Downloading: "Baixando", Completed: "Concluído",
		Error: "Erro", Cancelled: "Cancelado",
		EnterURL: "Insira uma URL", InvalidURL: "URL inválida",
		ErrorTitle: "Erro", UnknownFileSize: "Tamanho desconhecido",
		File: "Arquivo", Done: "Pronto",
		StatusBarTotal: "Total:", StatusBarActive: "Ativo:", StatusBarSpeed: "Veloc.:",
		ParallelConns: "conexões paralelas", SingleConnection: "conexão única",
	},
	ES: {
		WindowTitle: "Gestor de Descargas Rápido", NewDownload: "Nueva Descarga",
		URL: "URL:", Connections: "Conex.:", Download: "Descargar",
		CancelSelected: "Cancelar Selección", ClearCompleted: "Limpar Completados",
		Filename: "Archivo", Size: "Tamaño", Progress: "Progreso", Speed: "Veloc.",
		Time: "Tiempo", ETA: "Restante", Status: "Estado", Conns: "Conex.",
		Total: "Total:", Active: "Activo:", Language: "Idioma",
		Queued: "En cola", Downloading: "Descargando", Completed: "Completado",
		Error: "Error", Cancelled: "Cancelado",
		EnterURL: "Ingrese una URL", InvalidURL: "URL inválida",
		ErrorTitle: "Error", UnknownFileSize: "Tamaño desconocido",
		File: "Archivo", Done: "Hecho",
		StatusBarTotal: "Total:", StatusBarActive: "Activo:", StatusBarSpeed: "Veloc.:",
		ParallelConns: "conexiones paralelas", SingleConnection: "conexión única",
	},
	HI: {
		WindowTitle: "फास्ट डाउनलोड मैनेजर", NewDownload: "नया डाउनलोड",
		URL: "URL:", Connections: "कनेक्शन:", Download: "डाउनलोड",
		CancelSelected: "चयन रद्द करें", ClearCompleted: "पूर्ण हटाएं",
		Filename: "फ़ाइल नाम", Size: "आकार", Progress: "प्रगति", Speed: "गति",
		Time: "समय", ETA: "शेष", Status: "स्थिति", Conns: "कने.",
		Total: "कुल:", Active: "सक्रिय:", Language: "भाषा",
		Queued: "कतार में", Downloading: "डाउनलोड हो रहा", Completed: "पूर्ण",
		Error: "त्रुटि", Cancelled: "रद्द",
		EnterURL: "कृपया URL दर्ज करें", InvalidURL: "अमान्य URL",
		ErrorTitle: "त्रुटि", UnknownFileSize: "अज्ञात फ़ाइल आकार",
		File: "फ़ाइल", Done: "हो गया",
		StatusBarTotal: "कुल:", StatusBarActive: "सक्रिय:", StatusBarSpeed: "गति:",
		ParallelConns: "समानांतर कनेक्शन", SingleConnection: "एकल कनेक्शन",
	},
	CA: {
		WindowTitle: "Gestor de Descàrregues Ràpid", NewDownload: "Nova Descàrrega",
		URL: "URL:", Connections: "Connex.:", Download: "Descarregar",
		CancelSelected: "Cancelar Selecció", ClearCompleted: "Netejar Completats",
		Filename: "Nom del fitxer", Size: "Mida", Progress: "Progrés", Speed: "Veloc.",
		Time: "Temps", ETA: "Restant", Status: "Estat", Conns: "Conx.",
		Total: "Total:", Active: "Actiu:", Language: "Idioma",
		Queued: "En cua", Downloading: "Descarregant", Completed: "Completat",
		Error: "Error", Cancelled: "Cancel·lat",
		EnterURL: "Introduïu una URL", InvalidURL: "URL invàlida",
		ErrorTitle: "Error", UnknownFileSize: "Mida desconeguda",
		File: "Fitxer", Done: "Fet",
		StatusBarTotal: "Total:", StatusBarActive: "Actiu:", StatusBarSpeed: "Veloc.:",
		ParallelConns: "connexions paral·leles", SingleConnection: "connexió única",
	},
	RU: {
		WindowTitle: "Быстрый Менеджер Загрузок", NewDownload: "Новая загрузка",
		URL: "URL:", Connections: "Соед.:", Download: "Скачать",
		CancelSelected: "Отменить выбранное", ClearCompleted: "Очистить завершённые",
		Filename: "Имя файла", Size: "Размер", Progress: "Прогресс", Speed: "Скорость",
		Time: "Время", ETA: "Осталось", Status: "Статус", Conns: "Соед.",
		Total: "Всего:", Active: "Активно:", Language: "Язык",
		Queued: "В очереди", Downloading: "Загрузка", Completed: "Завершено",
		Error: "Ошибка", Cancelled: "Отменено",
		EnterURL: "Введите URL", InvalidURL: "Неверный URL",
		ErrorTitle: "Ошибка", UnknownFileSize: "Неизвестный размер",
		File: "Файл", Done: "Готово",
		StatusBarTotal: "Всего:", StatusBarActive: "Активно:", StatusBarSpeed: "Скорость:",
		ParallelConns: "параллельных соединений", SingleConnection: "одиночное соединение",
	},
	AZ: {
		WindowTitle: "Sürətli Yükləmə Meneceri", NewDownload: "Yeni Yükləmə",
		URL: "URL:", Connections: "Qoşulma:", Download: "Yüklə",
		CancelSelected: "Seçiləni Ləğv Et", ClearCompleted: "Bitənləri Təmizlə",
		Filename: "Fayl Adı", Size: "Ölçü", Progress: "Davam", Speed: "Sürət",
		Time: "Vaxt", ETA: "Qalan", Status: "Status", Conns: "Qoş.",
		Total: "Cəmi:", Active: "Aktiv:", Language: "Dil",
		Queued: "Növbədə", Downloading: "Yüklənir", Completed: "Tamam",
		Error: "Xəta", Cancelled: "Ləğv Edildi",
		EnterURL: "Zəhmət olmasa URL daxil edin", InvalidURL: "Yanlış URL",
		ErrorTitle: "Xəta", UnknownFileSize: "Fayl ölçüsü bilinmir",
		File: "Fayl", Done: "Hazır",
		StatusBarTotal: "Cəmi:", StatusBarActive: "Aktiv:", StatusBarSpeed: "Sürət:",
		ParallelConns: "paralel qoşulma", SingleConnection: "tək qoşulma",
	},
	KK: {
		WindowTitle: "Жылдам Жүктеу Менеджері", NewDownload: "Жаңа Жүктеу",
		URL: "URL:", Connections: "Қосылу:", Download: "Жүктеу",
		CancelSelected: "Таңдауды Болдырмау", ClearCompleted: "Аяқталғандарды Тазалау",
		Filename: "Файл аты", Size: "Өлшем", Progress: "Прогресс", Speed: "Жылд.",
		Time: "Уақыт", ETA: "Қалған", Status: "Күй", Conns: "Қос.",
		Total: "Барлығы:", Active: "Белсенді:", Language: "Тіл",
		Queued: "Кезекте", Downloading: "Жүктелуде", Completed: "Аяқталды",
		Error: "Қате", Cancelled: "Бас тартылды",
		EnterURL: "URL енгізіңіз", InvalidURL: "Жарамсыз URL",
		ErrorTitle: "Қате", UnknownFileSize: "Файл өлшемі белгісіз",
		File: "Файл", Done: "Дайын",
		StatusBarTotal: "Барлығы:", StatusBarActive: "Белсенді:", StatusBarSpeed: "Жылд.:",
		ParallelConns: "параллель қосылу", SingleConnection: "жалғыз қосылу",
	},
	UZ: {
		WindowTitle: "Tez Yuklab Oluvchi Menejer", NewDownload: "Yangi Yuklash",
		URL: "URL:", Connections: "Ulanish:", Download: "Yuklash",
		CancelSelected: "Tanlovni Bekor Qilish", ClearCompleted: "Tugallanganlarni Tozalash",
		Filename: "Fayl nomi", Size: "Hajm", Progress: "Jarayon", Speed: "Tezlik",
		Time: "Vaqt", ETA: "Qolgan", Status: "Holat", Conns: "Ulan.",
		Total: "Jami:", Active: "Faol:", Language: "Til",
		Queued: "Navbatda", Downloading: "Yuklanmoqda", Completed: "Tugallandi",
		Error: "Xato", Cancelled: "Bekor qilindi",
		EnterURL: "Iltimos URL kiriting", InvalidURL: "Noto'g'ri URL",
		ErrorTitle: "Xato", UnknownFileSize: "Fayl hajmi noma'lum",
		File: "Fayl", Done: "Tayyor",
		StatusBarTotal: "Jami:", StatusBarActive: "Faol:", StatusBarSpeed: "Tezlik:",
		ParallelConns: "parallel ulanish", SingleConnection: "yagona ulanish",
	},
}

func T(key string) string {
	p := langs[currentLang]
	switch key {
	case "WindowTitle": return p.WindowTitle
	case "NewDownload": return p.NewDownload
	case "URL": return p.URL
	case "Connections": return p.Connections
	case "Download": return p.Download
	case "CancelSelected": return p.CancelSelected
	case "ClearCompleted": return p.ClearCompleted
	case "Filename": return p.Filename
	case "Size": return p.Size
	case "Progress": return p.Progress
	case "Speed": return p.Speed
	case "Time": return p.Time
	case "ETA": return p.ETA
	case "Status": return p.Status
	case "Conns": return p.Conns
	case "Total": return p.Total
	case "Active": return p.Active
	case "Language": return p.Language
	case "Queued": return p.Queued
	case "Downloading": return p.Downloading
	case "Completed": return p.Completed
	case "Error": return p.Error
	case "Cancelled": return p.Cancelled
	case "EnterURL": return p.EnterURL
	case "InvalidURL": return p.InvalidURL
	case "ErrorTitle": return p.ErrorTitle
	case "UnknownFileSize": return p.UnknownFileSize
	case "File": return p.File
	case "Done": return p.Done
	case "StatusBarTotal": return p.StatusBarTotal
	case "StatusBarActive": return p.StatusBarActive
	case "StatusBarSpeed": return p.StatusBarSpeed
	}
	return key
}
