package models

import (
	"github.com/tebeka/selenium"
)

type QQZone struct {
	s *Selenium
}

func NewQQZone(s *Selenium) *QQZone {
	q := &QQZone{s: s}
	return q
}

func (self *QQZone) GET() {
	self.s.GET("https://i.qq.com/")
}

func (self *QQZone) Login(name, password string) {
	self.s.OnResponseCallback(func(wd selenium.WebDriver) {
		wd.SwitchFrame("login_frame")
		submit, _ := wd.FindElement(selenium.ByID, "switcher_plogin")
		submit.Click()
		u, _ := wd.FindElement(selenium.ByName, "u")
		u.Clear()
		u.SendKeys(name)

		p, _ := wd.FindElement(selenium.ByName, "p")
		p.Clear()
		p.SendKeys(password)

		wd.ExecuteScript("document.getElementById('login_button').parentNode.hidefocus=false;", nil)
		loginForm, _ := wd.FindElement(selenium.ByXPATH, `//*[@id="loginform"]/div[4]/a`)
		loginForm.Click()
		loginBtn, _ := wd.FindElement(selenium.ByID, "login_button")
		loginBtn.Click()
	})
}


