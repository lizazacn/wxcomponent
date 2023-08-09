package Global

import (
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/Struct"
	"sync"
	"time"
)

var EndTime time.Time
var Conf *Struct.Conf
var AccessToken string
var TokenLock sync.RWMutex
var QuestionTag string
var AnswerTag string
var MaxCallback int
