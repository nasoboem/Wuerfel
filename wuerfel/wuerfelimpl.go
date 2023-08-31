package wuerfel

import (. "würfel/zufallszahlen"
		"fmt"
		. "würfel/gfx"
		core "würfel/core"
		)

type data struct {
	seiten uint
	wert uint
	schatten uint
	wr,wg,wb uint8
	pr,pg,pb uint8
	hr,hg,hb uint8
	x,y,size uint16
	highlight bool
}

func New(seiten, schatten uint) *data {
	var w *data
	w = new(data)
	(*w).seiten = seiten
	(*w).schatten = schatten
	(*w).wr,(*w).wg,(*w).wb = 255,255,255
	(*w).pr,(*w).pg,(*w).pb = 0,0,0
	(*w).hr,(*w).hg,(*w).hb = 255,0,0
	//(*w).wert = 0         	//Alle auskommentierte Zeilen werden von Go automatisch mit Null initialisiert
	//(*w).x = 0				//Position und Größe wird über die Draw-Funktion gesteuert
	//(*w).y = 0
	//(*w).size = 0
	//(*w).highlight = false
	return w
}

//Geter-Funktionen
func (w *data) GibWert () uint {
	return (*w).wert
}

func (w *data) GibSeiten () uint {
	return (*w).seiten
}

func (w *data) GibWuerfelFarbe () (r,g,b uint8) {
	return (*w).wr,(*w).wg,(*w).wb
}

func (w *data) GibPunktFarbe () (r,g,b uint8) {
	return (*w).pr,(*w).pg,(*w).pb
}

func (w *data) GibPosition () (x,y uint16) {
	return (*w).x,(*w).y
}

func (w *data) GibGroesse () (size uint16) {
	return (*w).size
}

func (w *data) GibHighlightFarbe() (r,g,b uint8) {
	return (*w).hr,(*w).hg,(*w).hb
}

//Seter-Funktionen
func (w *data) SetzeWert (wert uint) { //kein Effekt bei Falscheingabe
	if wert <= (*w).seiten {
		(*w).wert = wert
	}
}

func (w *data) SetzeWertb (wert uint) { //bei Falscheingabe wird Wert maximal
	if wert <= (*w).seiten {
		(*w).wert = wert
	} else {
		(*w).wert = (*w).seiten
	}
}

func (w *data) SetzeWertc (wert uint) { //bei Falscheingabe bricht das Programm mit einer Panic ab
	if wert <= (*w).seiten {
		(*w).wert = wert
	} else {
		panic ("Gewünschter Wert größer als die Seitenzahl des Würfels!!!")
	}
}

func (w *data) SetzeSeiten (seiten uint) {
	(*w).seiten = seiten
	(*w).wert = 0
}

func (w *data) SetzeWuerfelFarbe (r,g,b uint8) {
	(*w).wr = r
	(*w).wg = g
	(*w).wb = b
}

func (w *data) SetzePunktFarbe (r,g,b uint8) {
	(*w).pr = r
	(*w).pg = g
	(*w).pb = b
}

func (w *data) SetzeHighlightFarbe (r,g,b uint8) {
	(*w).hr = r
	(*w).hg = g
	(*w).hb = b
}

func (w *data) Zuruecksetzen () {
	(*w).wert = 0
}

func (w *data) SetzeHighlight (highlight bool) {
	(*w).highlight = highlight
}

//Es gibt keine Seter-Funktion für die Position oder die Größe, da dies die Draw-Funktion mit ihrer Hilfsfunktion wasDrawnAs übernimmt!!!

func (w *data) String () string {
	var erg string
	//erg = erg + fmt.Sprint("[")
	erg = erg + fmt.Sprint((*w).wert)
	//erg = erg + fmt.Sprint("]")
	return erg
}


func (w *data) Wuerfeln () {
	Randomisieren()
	if (*w).seiten>0{
		(*w).wert = uint(Zufallszahl(1,int64((*w).seiten)))
	}
}

//Eine etwas aufwändigere Funktion, da wir hier alle programmierten Seitenanzahlen abdecken müssen, und diese alle einzeln, wie in der Draw-Funktion programmiert werden müssen.
func (w *data) PunktgehörtzumWuerfel (xp,yp uint16) bool {
	var xw,yw,size uint16
	xw = (*w).x
	yw = (*w).y
	size = (*w).size 
	switch (*w).seiten {
		case 0:
			switch (*w).schatten {
				case 2:
					if (int(xp)-int(xw))*(int(xp)-int(xw))+(int(yp)-int(yw))*(int(yp)-int(yw)) <= (int(size)) * (int(size)) {
						return true
					}
				case 4:
					//Volldreieck (x,y,x + size*2,y,x+size,y-2*size)
					if gehörtPunktzuDreieck(xp,yp,xw,yw,xw + size*2,yw,xw+size,yw-2*size) {
						return true
					}
					//Vollkreis (x,y,size/10)
					if (int(xp)-int(xw))*(int(xp)-int(xw))+(int(yp)-int(yw))*(int(yp)-int(yw)) <= int(size/10)*int(size/10) {
						return true
					}
					//Vollkreis (x + size*2,y,size/10)
					if (int(xp)-int(xw+size*2))*(int(xp)-int(xw+size*2))+(int(yp)-int(yw))*(int(yp)-int(yw))<= int(size/10)*int(size/10) {
						return true
					}
					//Vollkreis (x+size,y-2*size, size/10)
					if (int(xp)-int(xw+size))*(int(xp)-int(xw+size))+ (int(yp)-int(yw-2*size))*(int(yp)-int(yw-2*size))<= int(size/10)*int(size/10) {
						return true
					}
					//Volldreieck (x-size/10-1,y,x + size*2,y,x+size-size/10-1,y-2*size)
					if gehörtPunktzuDreieck(xp,yp,xw-size/10-1,yw,xw + size*2,yw,xw+size-size/10-1,yw-2*size) {
						return true
					}
					//Volldreieck (x,y,x + size*2+size/10+1,y,x+size+size/10+1,y-2*size)
					if gehörtPunktzuDreieck(xp,yp,xw,yw,xw + size*2+size/10+1,yw,xw+size+size/10+1,yw-2*size) {
						return true
					}
					//Volldreieck (x,y+size/10-1,x + size*2,y+size/10-1,x+size,y-2*size-size/10)
					if gehörtPunktzuDreieck(xp,yp,xw,yw+size/10-1,xw + size*2,yw+size/10-1,xw+size,yw-2*size-size/10) {
						return true
					}
				case 6:
					
					//Vollrechteck(x,y,size,size)
					if xp >= xw && xp <= xw + size && yp >= yw && yp <= yw + size {
						return true
					} 
					//Vollrechteck(x-size/10,y,size/10,size)
					if xp >= xw-size/10 && xp <= xw && yp >= yw && yp <= yw + size {
						return true
					}
					//Vollrechteck(x,y-size/10,size,size/10)
					if xp >= xw && xp <= xw + size && yp >= yw-size/10 && yp <= yw {
						return true
					}
					//Vollrechteck(x+size,y,size/10+1,size)
					if xp >= xw+size && xp <= xw+size+size/10+1 && yp >= yw && yp <= yw + size {
						return true
					}
					//Vollrechteck(x,y+size,size,size/10+1)
					if xp >= xw && xp <= xw + size && yp >= yw+size && yp <= yw+size+size/10+1 {
						return true
					}
					//Vollkreis(x,y,size/10) - hier wird der Satz des Pytagoras angewendet: Achtung!! bei den Variablen handelt es sich um uint16-Werte, für die Berechnung wird aber int benötigt!!! 
					if (int(xp)-int(xw))*(int(xp)-int(xw))+(int(yp)-int(yw))*(int(yp)-int(yw)) <= (int(size)/10) * (int(size)/10) {
						return true
					}
					//Vollkreis(x,y+size,size/10)
					if (int(xp)-int(xw))*(int(xp)-int(xw))+(int(yp)-int(yw+size))*(int(yp)-int(yw+size)) <= (int(size)/10) * (int(size)/10) {
						return true
					}
					//Vollkreis(x+size,y,size/10)
					if (int(xp)-int(xw+size))*(int(xp)-int(xw+size))+(int(yp)-int(yw))*(int(yp)-int(yw)) <= (int(size)/10) * (int(size)/10) {
						return true
					}
					//Vollkreis(x+size,y+size,size/10)
					if (int(xp)-int(xw+size))*(int(xp)-int(xw+size))+(int(yp)-int(yw+size))*(int(yp)-int(yw+size)) <= (int(size)/10) * (int(size)/10) {
						return true
					} 
				case 8:
					//Volldreieck (x,y,x + size*2,y,x+size,y-2*size)
					if gehörtPunktzuDreieck(xp,yp,xw,yw,xw + size*2,yw,xw+size,yw-2*size) {
						return true
					}
					//Vollkreis (x,y,size/10)
					if (int(xp)-int(xw))*(int(xp)-int(xw))+(int(yp)-int(yw))*(int(yp)-int(yw)) <= int(size/10)*int(size/10) {
						return true
					}
					//Vollkreis (x + size*2,y,size/10)
					if (int(xp)-int(xw+size*2))*(int(xp)-int(xw+size*2))+(int(yp)-int(yw))*(int(yp)-int(yw))<= int(size/10)*int(size/10) {
						return true
					}
					//Vollkreis (x+size,y-2*size, size/10)
					if (int(xp)-int(xw+size))*(int(xp)-int(xw+size))+ (int(yp)-int(yw-2*size))*(int(yp)-int(yw-2*size))<= int(size/10)*int(size/10) {
						return true
					}
					//Volldreieck (x-size/10-1,y,x + size*2,y,x+size-size/10-1,y-2*size)
					if gehörtPunktzuDreieck(xp,yp,xw-size/10-1,yw,xw + size*2,yw,xw+size-size/10-1,yw-2*size) {
						return true
					}
					//Volldreieck (x,y,x + size*2+size/10+1,y,x+size+size/10+1,y-2*size)
					if gehörtPunktzuDreieck(xp,yp,xw,yw,xw + size*2+size/10+1,yw,xw+size+size/10+1,yw-2*size) {
						return true
					}
					//Volldreieck (x,y+size/10-1,x + size*2,y+size/10-1,x+size,y-2*size-size/10)
					if gehörtPunktzuDreieck(xp,yp,xw,yw+size/10-1,xw + size*2,yw+size/10-1,xw+size,yw-2*size-size/10) {
						return true
					}
					//Volldreieck (x,y,x + size*2,y,x+size,y+size)
					if gehörtPunktzuDreieck(xp,yp,xw,yw,xw + size*2,yw,xw+size,yw+size) {
						return true
					}
					//Vollkreis   (x+size,y+size-size*2/100+size/2,size/10)
					if (int(xp)-int(xw+size))*(int(xp)-int(xw+size))+(int(yp)-int(yw+size-size*2/100+size/2))*(int(yp)-int(yw+size-size*2/100+size/2)) <= int(size/10)*int(size/10) {
						return true
					}
					//Volldreieck (x-size/10-1,y,x + size*2,y,x+size-size/10-1,y+size+size/2)
					if gehörtPunktzuDreieck(xp,yp,xw-size/10-1,yw,xw + size*2,yw,xw+size-size/10-1,yw+size+size/2) {
						return true
					}
					//Volldreieck (x,y,x+size*2+size/10+1,y,x+size+size/10+1,y+size+size/2)
					if gehörtPunktzuDreieck(xp,yp,xw,yw,xw+size*2+size/10+1,yw,xw+size+size/10+1,yw+size+size/2) {
						return true
					}
					//Volldreieck (x,y+size/10-1,x + size*2,y+size/10-1,x+size,y+size+size/2)
					if gehörtPunktzuDreieck(xp,yp,xw,yw+size/10-1,xw+size*2,yw+size/10-1,xw+size,yw+size+size/2) {
						return true
					}
				case 12:
					if (int(xp)-int(xw))*(int(xp)-int(xw))+(int(yp)-int(yw))*(int(yp)-int(yw)) <= (int(size*2/3)) * (int(size*2/3)) {
						return true
					}
				case 20:
					//Volldreieck(x,y+size*24/10,x+size*43/10,y,x+size*43/10,y+size*24/10)							//Dreieck 1
				if gehörtPunktzuDreieck(xp,yp,xw,yw+size*24/10,xw+size*43/10,yw,xw+size*43/10,yw+size*24/10) {
					return true
				}
				//Volldreieck(x+size*43/10,y,x+size*43/10,y+size*24/10,x+size*86/10,y+size*24/10)				//Dreieck 2
				if gehörtPunktzuDreieck(xp,yp,xw+size*43/10,yw,xw+size*43/10,yw+size*24/10,xw+size*86/10,yw+size*24/10) {
					return true
				}
				//Volldreieck(x,y+size*24/10,x+size*43/10,y+size*24/10,x+size*21/10,y+size*61/10)				//Dreisck 3
				if gehörtPunktzuDreieck(xp,yp,xw,yw+size*24/10,xw+size*43/10,yw+size*24/10,xw+size*21/10,yw+size*61/10) {
					return true
				}
				//Volldreieck(x+size*43/10,y+size*24/10,x+size*21/10,y+size*61/10,x+size*65/10,y+size*61/10)	//Dreieck 4
				if gehörtPunktzuDreieck(xp,yp,xw+size*43/10,yw+size*24/10,xw+size*21/10,yw+size*61/10,xw+size*65/10,yw+size*61/10) {
					return true
				}
				//Volldreieck(x+size*43/10,y+size*24/10,x+size*86/10,y+size*24/10,x+size*65/10,y+size*61/10)	//Dreieck 5
				if gehörtPunktzuDreieck(xp,yp,xw+size*43/10,yw+size*24/10,xw+size*86/10,yw+size*24/10,xw+size*65/10,yw+size*61/10) {
					return true
				}
				//Volldreieck(x,y+size*24/10,x+size*21/10,y+size*61/10,x,y+size*73/10)							//Dreieck 6
				if gehörtPunktzuDreieck(xp,yp,xw,yw+size*24/10,xw+size*21/10,yw+size*61/10,xw,yw+size*73/10) {
					return true
				}
				//Volldreieck(x+size*86/10,y+size*24/10,x+size*86/10,y+size*73/10,x+size*65/10,y+size*61/10)	//Dreieck 7
				if gehörtPunktzuDreieck(xp,yp,xw+size*86/10,yw+size*24/10,xw+size*86/10,yw+size*73/10,xw+size*65/10,yw+size*61/10) {
					return true
				}
				//Volldreieck(x,y+size*73/10,x+size*21/10,y+size*61/10,x+size*43/10,y+size*98/10)				//Dreieck 8
				if gehörtPunktzuDreieck(xp,yp,xw,yw+size*73/10,xw+size*21/10,yw+size*61/10,xw+size*43/10,yw+size*98/10) {
					return true
				}
				//Volldreieck(x+size*21/10,y+size*61/10,x+size*65/10,y+size*61/10,x+size*43/10,y+size*98/10)	//Dreieck 9
				if gehörtPunktzuDreieck(xp,yp,xw+size*21/10,yw+size*61/10,xw+size*65/10,yw+size*61/10,xw+size*43/10,yw+size*98/10) {
					return true
				}
				//Volldreieck(x+size*65/10,y+size*61/10,x+size*86/10,y+size*73/10,x+size*43/10,y+size*98/10)	//Dreieck 10
				if gehörtPunktzuDreieck(xp,yp,xw+size*65/10,yw+size*61/10,xw+size*86/10,yw+size*73/10,xw+size*43/10,yw+size*98/10) {
					return true
				}
				
			}
		case 2:
			if (int(xp)-int(xw))*(int(xp)-int(xw))+(int(yp)-int(yw))*(int(yp)-int(yw)) <= (int(size)) * (int(size)) {
						return true
			}
		case 4:
			//Volldreieck (x,y,x + size*2,y,x+size,y-2*size)
			if gehörtPunktzuDreieck(xp,yp,xw,yw,xw + size*2,yw,xw+size,yw-2*size) {
				return true
			}
			//Vollkreis (x,y,size/10)
			if (int(xp)-int(xw))*(int(xp)-int(xw))+(int(yp)-int(yw))*(int(yp)-int(yw)) <= int(size/10)*int(size/10) {
				return true
			}
			//Vollkreis (x + size*2,y,size/10)
			if (int(xp)-int(xw+size*2))*(int(xp)-int(xw+size*2))+(int(yp)-int(yw))*(int(yp)-int(yw))<= int(size/10)*int(size/10) {
				return true
			}
			//Vollkreis (x+size,y-2*size, size/10)
			if (int(xp)-int(xw+size))*(int(xp)-int(xw+size))+ (int(yp)-int(yw-2*size))*(int(yp)-int(yw-2*size))<= int(size/10)*int(size/10) {
				return true
			}
			//Volldreieck (x-size/10-1,y,x + size*2,y,x+size-size/10-1,y-2*size)
			if gehörtPunktzuDreieck(xp,yp,xw-size/10-1,yw,xw + size*2,yw,xw+size-size/10-1,yw-2*size) {
				return true
			}
			//Volldreieck (x,y,x + size*2+size/10+1,y,x+size+size/10+1,y-2*size)
			if gehörtPunktzuDreieck(xp,yp,xw,yw,xw + size*2+size/10+1,yw,xw+size+size/10+1,yw-2*size) {
				return true
			}
			//Volldreieck (x,y+size/10-1,x + size*2,y+size/10-1,x+size,y-2*size-size/10)
			if gehörtPunktzuDreieck(xp,yp,xw,yw+size/10-1,xw + size*2,yw+size/10-1,xw+size,yw-2*size-size/10) {
				return true
			}
		case 6:
			//Vollrechteck(x,y,size,size)
			if xp >= xw && xp <= xw + size && yp >= yw && yp <= yw + size {
				return true
			} 
			//Vollrechteck(x-size/10,y,size/10,size)
			if xp >= xw-size/10 && xp <= xw && yp >= yw && yp <= yw + size {
				return true
			}
			//Vollrechteck(x,y-size/10,size,size/10)
			if xp >= xw && xp <= xw + size && yp >= yw-size/10 && yp <= yw {
				return true
			}
			//Vollrechteck(x+size,y,size/10+1,size)
			if xp >= xw+size && xp <= xw+size+size/10+1 && yp >= yw && yp <= yw + size {
				return true
			}
			//Vollrechteck(x,y+size,size,size/10+1)
			if xp >= xw && xp <= xw + size && yp >= yw+size && yp <= yw+size+size/10+1 {
				return true
			}
			//Vollkreis(x,y,size/10) - hier wird der Satz des Pytagoras angewendet: Achtung!! bei den Variablen handelt es sich um uint16-Werte, für die Berechnung wird aber int benötigt!!! 
			if (int(xp)-int(xw))*(int(xp)-int(xw))+(int(yp)-int(yw))*(int(yp)-int(yw)) <= (int(size)/10) * (int(size)/10) {
				return true
			}
			//Vollkreis(x,y+size,size/10)
			if (int(xp)-int(xw))*(int(xp)-int(xw))+(int(yp)-int(yw+size))*(int(yp)-int(yw+size)) <= (int(size)/10) * (int(size)/10) {
				return true
			}
			//Vollkreis(x+size,y,size/10)
			if (int(xp)-int(xw+size))*(int(xp)-int(xw+size))+(int(yp)-int(yw))*(int(yp)-int(yw)) <= (int(size)/10) * (int(size)/10) {
				return true
			}
			//Vollkreis(x+size,y+size,size/10)
			if (int(xp)-int(xw+size))*(int(xp)-int(xw+size))+(int(yp)-int(yw+size))*(int(yp)-int(yw+size)) <= (int(size)/10) * (int(size)/10) {
				return true
			} 
		case 8:
			//Volldreieck (x,y,x + size*2,y,x+size,y-2*size)
			if gehörtPunktzuDreieck(xp,yp,xw,yw,xw + size*2,yw,xw+size,yw-2*size) {
				return true
			}
			//Vollkreis (x,y,size/10)
			if (int(xp)-int(xw))*(int(xp)-int(xw))+(int(yp)-int(yw))*(int(yp)-int(yw)) <= int(size/10)*int(size/10) {
				return true
			}
			//Vollkreis (x + size*2,y,size/10)
			if (int(xp)-int(xw+size*2))*(int(xp)-int(xw+size*2))+(int(yp)-int(yw))*(int(yp)-int(yw))<= int(size/10)*int(size/10) {
				return true
			}
			//Vollkreis (x+size,y-2*size, size/10)
			if (int(xp)-int(xw+size))*(int(xp)-int(xw+size))+ (int(yp)-int(yw-2*size))*(int(yp)-int(yw-2*size))<= int(size/10)*int(size/10) {
				return true
			}
			//Volldreieck (x-size/10-1,y,x + size*2,y,x+size-size/10-1,y-2*size)
			if gehörtPunktzuDreieck(xp,yp,xw-size/10-1,yw,xw + size*2,yw,xw+size-size/10-1,yw-2*size) {
				return true
			}
			//Volldreieck (x,y,x + size*2+size/10+1,y,x+size+size/10+1,y-2*size)
			if gehörtPunktzuDreieck(xp,yp,xw,yw,xw + size*2+size/10+1,yw,xw+size+size/10+1,yw-2*size) {
				return true
			}
			//Volldreieck (x,y+size/10-1,x + size*2,y+size/10-1,x+size,y-2*size-size/10)
			if gehörtPunktzuDreieck(xp,yp,xw,yw+size/10-1,xw + size*2,yw+size/10-1,xw+size,yw-2*size-size/10) {
				return true
			}
			
			//Volldreieck (x,y,x + size*2,y,x+size,y+size)
			if gehörtPunktzuDreieck(xp,yp,xw,yw,xw + size*2,yw,xw+size,yw+size) {
				return true
			}
			//Vollkreis   (x+size,y+size-size*2/100+size/2,size/10)
			if (int(xp)-int(xw+size))*(int(xp)-int(xw+size))+(int(yp)-int(yw+size-size*2/100+size/2))*(int(yp)-int(yw+size-size*2/100+size/2)) <= int(size/10)*int(size/10) {
				return true
			}
			//Volldreieck (x-size/10-1,y,x + size*2,y,x+size-size/10-1,y+size+size/2)
			if gehörtPunktzuDreieck(xp,yp,xw-size/10-1,yw,xw + size*2,yw,xw+size-size/10-1,yw+size+size/2) {
				return true
			}
			//Volldreieck (x,y,x+size*2+size/10+1,y,x+size+size/10+1,y+size+size/2)
			if gehörtPunktzuDreieck(xp,yp,xw,yw,xw+size*2+size/10+1,yw,xw+size+size/10+1,yw+size+size/2) {
				return true
			}
			//Volldreieck (x,y+size/10-1,x + size*2,y+size/10-1,x+size,y+size+size/2)
			if gehörtPunktzuDreieck(xp,yp,xw,yw+size/10-1,xw+size*2,yw+size/10-1,xw+size,yw+size+size/2) {
				return true
			}
			
		case 12:
			if (int(xp)-int(xw))*(int(xp)-int(xw))+(int(yp)-int(yw))*(int(yp)-int(yw)) <= (int(size*2/3)) * (int(size*2/3)) {
				return true
			}
		case 20:
				//Volldreieck(x,y+size*24/10,x+size*43/10,y,x+size*43/10,y+size*24/10)							//Dreieck 1
				if gehörtPunktzuDreieck(xp,yp,xw,yw+size*24/10,xw+size*43/10,yw,xw+size*43/10,yw+size*24/10) {
					return true
				}
				//Volldreieck(x+size*43/10,y,x+size*43/10,y+size*24/10,x+size*86/10,y+size*24/10)				//Dreieck 2
				if gehörtPunktzuDreieck(xp,yp,xw+size*43/10,yw,xw+size*43/10,yw+size*24/10,xw+size*86/10,yw+size*24/10) {
					return true
				}
				//Volldreieck(x,y+size*24/10,x+size*43/10,y+size*24/10,x+size*21/10,y+size*61/10)				//Dreisck 3
				if gehörtPunktzuDreieck(xp,yp,xw,yw+size*24/10,xw+size*43/10,yw+size*24/10,xw+size*21/10,yw+size*61/10) {
					return true
				}
				//Volldreieck(x+size*43/10,y+size*24/10,x+size*21/10,y+size*61/10,x+size*65/10,y+size*61/10)	//Dreieck 4
				if gehörtPunktzuDreieck(xp,yp,xw+size*43/10,yw+size*24/10,xw+size*21/10,yw+size*61/10,xw+size*65/10,yw+size*61/10) {
					return true
				}
				//Volldreieck(x+size*43/10,y+size*24/10,x+size*86/10,y+size*24/10,x+size*65/10,y+size*61/10)	//Dreieck 5
				if gehörtPunktzuDreieck(xp,yp,xw+size*43/10,yw+size*24/10,xw+size*86/10,yw+size*24/10,xw+size*65/10,yw+size*61/10) {
					return true
				}
				//Volldreieck(x,y+size*24/10,x+size*21/10,y+size*61/10,x,y+size*73/10)							//Dreieck 6
				if gehörtPunktzuDreieck(xp,yp,xw,yw+size*24/10,xw+size*21/10,yw+size*61/10,xw,yw+size*73/10) {
					return true
				}
				//Volldreieck(x+size*86/10,y+size*24/10,x+size*86/10,y+size*73/10,x+size*65/10,y+size*61/10)	//Dreieck 7
				if gehörtPunktzuDreieck(xp,yp,xw+size*86/10,yw+size*24/10,xw+size*86/10,yw+size*73/10,xw+size*65/10,yw+size*61/10) {
					return true
				}
				//Volldreieck(x,y+size*73/10,x+size*21/10,y+size*61/10,x+size*43/10,y+size*98/10)				//Dreieck 8
				if gehörtPunktzuDreieck(xp,yp,xw,yw+size*73/10,xw+size*21/10,yw+size*61/10,xw+size*43/10,yw+size*98/10) {
					return true
				}
				//Volldreieck(x+size*21/10,y+size*61/10,x+size*65/10,y+size*61/10,x+size*43/10,y+size*98/10)	//Dreieck 9
				if gehörtPunktzuDreieck(xp,yp,xw+size*21/10,yw+size*61/10,xw+size*65/10,yw+size*61/10,xw+size*43/10,yw+size*98/10) {
					return true
				}
				//Volldreieck(x+size*65/10,y+size*61/10,x+size*86/10,y+size*73/10,x+size*43/10,y+size*98/10)	//Dreieck 10
				if gehörtPunktzuDreieck(xp,yp,xw+size*65/10,yw+size*61/10,xw+size*86/10,yw+size*73/10,xw+size*43/10,yw+size*98/10) {
					return true
				}
	}
	return false
}

//Hilfsfunktion zur Detektirung von Dreiecken

func gehörtPunktzuDreieck (xp,yp,x1,y1,x2,y2,x3,y3 uint16) bool {
	var p1,p2,p3 [2]uint16							//Eckpunkte des Dreiecks
	p1[0] = x1										//Umwandlung der x- und Y-Kooerdinaten in einen Punkt ([2]uint16)- Dreieck
	p1[1] = y1
	p2[0] = x2
	p2[1] = y2
	p3[0] = x3
	p3[1] = y3
	var x [2]uint16									//Umwandlung der x- und Y-Kooerdinaten in einen Punkt ([2]uint16) - Zu untersuchender Punkt
	x[0] = xp
	x[1] = yp
	return selbeSeite(p1,p2,p3,x) && selbeSeite(p1,p3,p2,x) && selbeSeite(p2,p3,p1,x)  //Test ob der Punkt auf der Innenseite aller Dreiseckkanten ist.
}

func selbeSeite (a,b,c,p [2]uint16) bool {
	var cp1,cp2 core.Vector
	var v1, v2, v3 core.Vector
	v1.X = float64(int(b[0])-int(a[0]))
	v1.Y = float64(int(b[1])-int(a[1]))
	v1.Z = 0
	v2.X = float64(int(c[0])-int(a[0]))
	v2.Y = float64(int(c[1])-int(a[1]))
	v2.Z = 0
	v3.X = float64(int(p[0])-int(a[0]))
	v3.Y = float64(int(p[1])-int(a[1]))
	v3.Z = 0
    cp1 = v1.Cross(v2)
    cp2 = v1.Cross(v3)
    return cp1.Dot(cp2) >= 0
}


//Hilfsfunktionen zur Veränderung von Größe und Position, mit Anpassung der Größe & Position

func (w *data) wasDrawnAs (x,y,size uint16) {
	var wtype uint
	if (*w).seiten==0{
		wtype = (*w).schatten
	} else {
		wtype = (*w).seiten
	} 
	switch wtype {
		case 2:
			size = size*6/5
			y = y + size*13/15
			x = x + size
		case 4:
			y = y + size*2
		case 6:
			size = size*2
		case 8:
			y = y + size*2
		case 12:
			size = size*16/8
			y = y + size*9/15
			x = x + size/2
		case 20:
			size = size*65/160
			y = y - size/2
	}
	(*w).x = x
	(*w).y = y
	(*w).size = size
} 


func (w *data) Draw (x,y,size uint16) {
	var zeilen,spalten uint16			//Fragt die Zeilen und Spalten des Grafikfensters ab spalten=pixel in x-Achse/Breite, zeilen = pixel in y-Achse/Höhe
	if FensterOffen() {
		zeilen = Grafikzeilen()
		spalten = Grafikspalten()
	}
	w.wasDrawnAs(x,y,size)				//Schreibt die Werte von x,y und size in den Würfel; s. Hilfsfunktion
	size = (*w).size
	x = (*w).x
	y = (*w).y
	var wr,wg,wb,pr,pg,pb,hr,hg,hb uint8
	wr,wg,wb = w.GibWuerfelFarbe()
	pr,pg,pb = w.GibPunktFarbe()
	hr,hg,hb = w.GibHighlightFarbe()
	switch (*w).seiten {
		case 0:
		
			switch (*w).schatten {
				case 2:
					if (*w).highlight {
						Stiftfarbe(hr,hg,hb)
						Vollkreis(x,y,size+5)
					}
				
					Stiftfarbe(120,120,120)
					Vollkreis(x,y,size)
				case 4:
					if (*w).highlight {
						Stiftfarbe(hr,hg,hb)
						Volldreieck (x,y,x + size*2,y,x+size,y-2*size)
						Vollkreis (x,y,size/10+5)
						Vollkreis (x + size*2,y,size/10+5)
						Vollkreis (x+size,y-2*size, size/10+5)
						Volldreieck (x-size/10-5,y-5,x + size*2,y,x+size-size/10-5,y-2*size-5)
						Volldreieck (x,y,x + size*2+size/10+5-1,y-5,x+size+size/10+5,y-2*size-5+1)
						Volldreieck (x,y+size/10-1+5,x + size*2,y+size/10-1+5,x+size,y-2*size-size/10)
					}
					Stiftfarbe(120,120,120)
					Volldreieck (x,y,x + size*2,y,x+size,y-2*size)
					Vollkreis (x,y,size/10)
					Vollkreis (x + size*2,y,size/10)
					Vollkreis (x+size,y-2*size, size/10)
					Volldreieck (x-size/10-1,y,x + size*2,y,x+size-size/10-1,y-2*size)
					Volldreieck (x,y,x + size*2+size/10+1,y,x+size+size/10+1,y-2*size)
					Volldreieck (x,y+size/10-1,x + size*2,y+size/10-1,x+size,y-2*size-size/10)
					
				case 6:
					if (*w).highlight {
						Stiftfarbe(hr,hg,hb)
						Vollrechteck(x-5,y-5,size+10,size+10)
						Vollrechteck(x-5-(size+10)/10,y-5,(size+10)/10,size+10)
						Vollrechteck(x-5,y-5 -(size+10)/10,size+10,(size+10)/10)
						Vollrechteck(x-5+size+10,y-5,(size+10)/10+1,size+10)
						Vollrechteck(x-5,y-5+size+10,size+10,(size+10)/10+1)
						Vollkreis(x-5,y-5,(size+10)/10)
						Vollkreis(x-5,y-5+size+10,(size+10)/10)
						Vollkreis(x-5+size+10,y-5,(size+10)/10)
						Vollkreis(x-5+size+10,y-5+size+10,(size+10)/10)
					}
					Stiftfarbe(120,120,120)
					Vollrechteck(x,y,size,size)
					Vollrechteck(x-size/10,y,size/10,size)
					Vollrechteck(x,y-size/10,size,size/10)
					Vollrechteck(x+size,y,size/10+1,size)
					Vollrechteck(x,y+size,size,size/10+1)
					Vollkreis(x,y,size/10)
					Vollkreis(x,y+size,size/10)
					Vollkreis(x+size,y,size/10)
					Vollkreis(x+size,y+size,size/10)
				case 8:
					if (*w).highlight {
						Stiftfarbe(hr,hg,hb)
						Volldreieck (x,y,x + size*2,y,x+size,y-2*size)
						Vollkreis (x,y,size/10+6)
						Vollkreis (x + size*2,y,size/10+6)
						Vollkreis (x+size,y-2*size, size/10+5)
						Volldreieck (x-size/10-5,y-5,x + size*2,y,x+size-size/10-5,y-2*size-5)
						Volldreieck (x,y,x + size*2+size/10+5-1,y-5,x+size+size/10+5,y-2*size-5+1)
						Vollkreis   (x+size,y+size-size*2/100+size/2,size/10+5)
						Volldreieck (x-size/10-6-1,y,x + size*2,y,x+size-size/10-6,y+size+size/2)
						Volldreieck (x,y,x+size*2+size/10+6+1,y,x+size+size/10+6,y+size+size/2)				
					}
		
					Stiftfarbe(120,120,120)
					Volldreieck (x,y,x + size*2,y,x+size,y+size)
					Vollkreis   (x+size,y+size-size*2/100+size/2,size/10)
					Volldreieck (x-size/10-1,y,x + size*2,y,x+size-size/10-1,y+size+size/2)
					Volldreieck (x,y,x+size*2+size/10+1,y,x+size+size/10+1,y+size+size/2)
					Volldreieck (x,y+size/10-1,x + size*2,y+size/10-1,x+size,y+size+size/2)
					Volldreieck (x,y,x + size*2,y,x+size,y-2*size)
					Vollkreis (x,y,size/10)
					Vollkreis (x + size*2,y,size/10)
					Vollkreis (x+size,y-2*size, size/10)
					Volldreieck (x-size/10-1,y,x + size*2,y,x+size-size/10-1,y-2*size)
					Volldreieck (x,y,x + size*2+size/10+1,y,x+size+size/10+1,y-2*size)
					Volldreieck (x,y+size/10-1,x + size*2,y+size/10-1,x+size,y-2*size-size/10)
				
				case 12:
					if (*w).highlight {
						Stiftfarbe(hr,hg,hb)
						Vollkreis(x,y,size*2/3+5)
					}
				
					Stiftfarbe(120,120,120)
					Vollkreis(x,y,size*2/3)
				case 20:
					if (*w).highlight {
						Stiftfarbe(hr,hg,hb)
						Volldreieck(x,y+size*24/10,x+size*43/10,y,x+size*43/10,y-7)											//Dreieck 1
						Volldreieck(x,y+size*24/10,x+size*43/10,y-7,x-5,y+size*24/10-5)										//Dreieck 2
						Volldreieck(x,y+size*24/10,x,y+size*73/10,x-5,y+size*24/10-5)										//Dreisck 3
						Volldreieck(x,y+size*73/10,x-5,y+size*24/10-5,x-5,y+size*73/10+5)									//Dreieck 4
						Volldreieck(x,y+size*73/10,x+size*43/10,y+size*98/10,x+size*43/10,y+size*98/10+7)					//Dreieck 5
						Volldreieck(x,y+size*73/10,x-5,y+size*73/10+5,x+size*43/10,y+size*98/10+7)							//Dreieck 6
						Volldreieck(x+size*43/10,y+size*98/10,x+size*43/10,y+size*98/10+7,x+size*86/10+5,y+size*73/10+5)	//Dreieck 7
						Volldreieck(x+size*43/10,y+size*98/10,x+size*86/10,y+size*73/10,x+size*86/10+5,y+size*73/10+5)		//Dreieck 8
						Volldreieck(x+size*86/10,y+size*73/10,x+size*86/10,y+size*24/10,x+size*86/10+5,y+size*24/10-5)		//Dreieck 9
						Volldreieck(x+size*86/10,y+size*73/10,x+size*86/10+5,y+size*73/10+5,x+size*86/10+5,y+size*24/10-5)	//Dreieck 10
						Volldreieck(x+size*86/10,y+size*24/10,x+size*86/10+5,y+size*24/10-5,x+size*43/10,y-7)				//Dreieck 11
						Volldreieck(x+size*43/10,y,x+size*86/10,y+size*24/10,x+size*43/10,y-7)								//Dreieck 12
					}
				
					Stiftfarbe(120,120,120)
					Volldreieck(x,y+size*24/10,x+size*43/10,y,x+size*43/10,y+size*24/10)						//Dreieck 1
					Volldreieck(x+size*43/10,y,x+size*43/10,y+size*24/10,x+size*86/10,y+size*24/10)				//Dreieck 2
					Volldreieck(x,y+size*24/10,x+size*43/10,y+size*24/10,x+size*21/10,y+size*61/10)				//Dreisck 3
					Volldreieck(x+size*43/10,y+size*24/10,x+size*21/10,y+size*61/10,x+size*65/10,y+size*61/10)	//Dreieck 4
					Volldreieck(x+size*43/10,y+size*24/10,x+size*86/10,y+size*24/10,x+size*65/10,y+size*61/10)	//Dreieck 5
					Volldreieck(x,y+size*24/10,x+size*21/10,y+size*61/10,x,y+size*73/10)						//Dreieck 6
					Volldreieck(x+size*86/10,y+size*24/10,x+size*86/10,y+size*73/10,x+size*65/10,y+size*61/10)	//Dreieck 7
					Volldreieck(x,y+size*73/10,x+size*21/10,y+size*61/10,x+size*43/10,y+size*98/10)				//Dreieck 8
					Volldreieck(x+size*21/10,y+size*61/10,x+size*65/10,y+size*61/10,x+size*43/10,y+size*98/10)	//Dreieck 9
					Volldreieck(x+size*65/10,y+size*61/10,x+size*86/10,y+size*73/10,x+size*43/10,y+size*98/10)	//Dreieck 10
					
				default: 
					panic("Sie haben ein Schatten eines Würfels gewählt, der nicht von dem Würfelpaket gezeichnet werden kann!! Wählen Sie New(0,2), New(0,4), New(0,6), New(0,8), New(0,12) oder New(0,20) oder ihre Invertierung!!!") 

			}
		
		case 2:		//Münze
			if (*w).highlight {
				Stiftfarbe(hr,hg,hb)
				Vollkreis(x,y,size+5)
			}
		
			//Münzkörper
			Stiftfarbe(wr,wg,wb)
			Vollkreis(x,y,size)
			Stiftfarbe(pr,pg,pb)
			Kreis (x,y,size)
			Kreis (x,y,size*17/18)
			
			//Wert
			switch (*w).wert {
				case 1:
					Vollkreis(x,y-size/18,size*10/18)
					Stiftfarbe(wr,wg,wb)
					Vollkreis(x,y-size/18,size*9/18)
					Stiftfarbe(pr,pg,pb)
					Vollkreis(x-(size*9/18)/2,y-(size*9/18)/2,size/12)
					Vollkreis(x+(size*9/18)/2,y-(size*9/18)/2,size/12)
					Vollkreissektor (x,y,(size*9/18)/2,180,360)
					Vollrechteck(x-size/20,y+(size*9/18),size/10,size/5)
					
				case 2:
					Vollrechteck (x-size/20,y-size/2,size/10,size-size/10)
					Vollrechteck (x-size/4,y+size/2-size/10,size/2,size/10)
					Vollrechteck (x-size/4,y-size/2,size/2,size/10)
			}
		case 4:		//Tetraeder-Würfel
			if (*w).highlight {
				Stiftfarbe(hr,hg,hb)
				Volldreieck (x,y,x + size*2,y,x+size,y-2*size)
				Vollkreis (x,y,size/10+5)
				Vollkreis (x + size*2,y,size/10+5)
				Vollkreis (x+size,y-2*size, size/10+5)
				Volldreieck (x-size/10-5,y-5,x + size*2,y,x+size-size/10-5,y-2*size-5)
				Volldreieck (x,y,x + size*2+size/10+5-1,y-5,x+size+size/10+5,y-2*size-5+1)
				Volldreieck (x,y+size/10-1+5,x + size*2,y+size/10-1+5,x+size,y-2*size-size/10)
			}
			Stiftfarbe(wr,wg,wb)
			Volldreieck (x,y,x + size*2,y,x+size,y-2*size)
			Vollkreis (x,y,size/10)
			Vollkreis (x + size*2,y,size/10)
			Vollkreis (x+size,y-2*size, size/10)
			Volldreieck (x-size/10-1,y,x + size*2,y,x+size-size/10-1,y-2*size)
			Volldreieck (x,y,x + size*2+size/10+1,y,x+size+size/10+1,y-2*size)
			Volldreieck (x,y+size/10-1,x + size*2,y+size/10-1,x+size,y-2*size-size/10)
			
			
			Stiftfarbe(pr,pg,pb)
			switch (*w).wert {
				case 1:
					Vollrechteck (x+size-size/20,y-size*11/10,size/10,size*4/5)
				case 2.:
					Vollrechteck (x+size*4/5,y-size*2/5,size*4/10,size/10)
					Vollrechteck (x+size*4/5,y-size*11/10,size*2/5,size/10)
					Vollrechteck (x+size*4/5,y-size*37/50,size*2/5,size/10)
					Vollrechteck (x+size*6/5,y-size,size/10,size*13/50)
					Vollrechteck (x+size*7/10,y-size*32/50,size/10,size*13/50)
				case 3:
					Vollrechteck (x+size*4/5,y-size*11/10,size*2/5,size/10)
					Vollrechteck (x+size*4/5,y-size*37/50,size*2/5,size/10)
					Vollrechteck (x+size*4/5,y-size*2/5,size*2/5,size/10)
					Vollrechteck (x+size*6/5,y-size,size/10,size*13/50)
					Vollrechteck (x+size*6/5,y-size*32/50,size/10,size*12/50)
				case 4:
					Vollrechteck (x+size*4/5,y-size*11/10,size/10,size*3/5)
					Vollrechteck (x+size*4/5,y-size*3/5,size*27/50,size/10)
					Vollrechteck (x+size*11/10,y-size*4/5,size/10,size*23/50)
			}
		
		case 6:		//normaler 6-seitiger Würfel
			if (*w).highlight {
				Stiftfarbe(hr,hg,hb)
				Vollrechteck(x-5,y-5,size+10,size+10)
				Vollrechteck(x-5-(size+10)/10,y-5,(size+10)/10,size+10)
				Vollrechteck(x-5,y-5 -(size+10)/10,size+10,(size+10)/10)
				Vollrechteck(x-5+size+10,y-5,(size+10)/10+1,size+10)
				Vollrechteck(x-5,y-5+size+10,size+10,(size+10)/10+1)
				Vollkreis(x-5,y-5,(size+10)/10)
				Vollkreis(x-5,y-5+size+10,(size+10)/10)
				Vollkreis(x-5+size+10,y-5,(size+10)/10)
				Vollkreis(x-5+size+10,y-5+size+10,(size+10)/10)
			}
				Stiftfarbe(pr,pg,pb)
			if x+size+size/10<spalten && y+size+size/10<zeilen && x-size/10>=1 && y-size/10>=1 {
				Stiftfarbe(wr,wg,wb)
				Vollrechteck(x,y,size,size)
				Vollrechteck(x-size/10,y,size/10,size)
				Vollrechteck(x,y-size/10,size,size/10)
				Vollrechteck(x+size,y,size/10+1,size)
				Vollrechteck(x,y+size,size,size/10+1)
				Vollkreis(x,y,size/10)
				Vollkreis(x,y+size,size/10)
				Vollkreis(x+size,y,size/10)
				Vollkreis(x+size,y+size,size/10)
				Stiftfarbe(pr,pg,pb)
				switch (*w).wert {
					case 1:
						Vollkreis(x+size/2,y+size/2,size/10)
					case 2:
						Vollkreis(x+size/5,y+size/5,size/10)
						Vollkreis(x+size*4/5,y+size*4/5,size/10)
					case 3:
						Vollkreis(x+size/5,y+size/5,size/10)
						Vollkreis(x+size/2,y+size/2,size/10)
						Vollkreis(x+size*4/5,y+size*4/5,size/10)
					case 4:
						Vollkreis(x+size/5,y+size/5,size/10)
						Vollkreis(x+size*4/5,y+size*4/5,size/10)
						Vollkreis(x+size/5,y+size*4/5,size/10)
						Vollkreis(x+size*4/5,y+size/5,size/10)
					case 5:
						Vollkreis(x+size/5,y+size/5,size/10)
						Vollkreis(x+size/2,y+size/2,size/10)
						Vollkreis(x+size*4/5,y+size*4/5,size/10)
						Vollkreis(x+size*4/5,y+size/5,size/10)
						Vollkreis(x+size/5,y+size*4/5,size/10)
					case 6:
						Vollkreis(x+size/5,y+size/5,size/10)
						Vollkreis(x+size/5,y+size/2,size/10)
						Vollkreis(x+size*4/5,y+size/2,size/10)
						Vollkreis(x+size*4/5,y+size*4/5,size/10)
						Vollkreis(x+size*4/5,y+size/5,size/10)
						Vollkreis(x+size/5,y+size*4/5,size/10)	
				}
			}

		case 8:		//Oktaeder
			if (*w).highlight {
				Stiftfarbe(hr,hg,hb)
				Volldreieck (x,y,x + size*2,y,x+size,y-2*size)
				Vollkreis (x,y,size/10+6)
				Vollkreis (x + size*2,y,size/10+6)
				Vollkreis (x+size,y-2*size, size/10+5)
				Volldreieck (x-size/10-5,y-5,x + size*2,y,x+size-size/10-5,y-2*size-5)
				Volldreieck (x,y,x + size*2+size/10+5-1,y-5,x+size+size/10+5,y-2*size-5+1)
				Vollkreis   (x+size,y+size-size*2/100+size/2,size/10+5)
				Volldreieck (x-size/10-6-1,y,x + size*2,y,x+size-size/10-6,y+size+size/2)
				Volldreieck (x,y,x+size*2+size/10+6+1,y,x+size+size/10+6,y+size+size/2)				
			}
		
			Stiftfarbe(130,130,130) //Schatten des unteren Teils
			Volldreieck (x,y,x + size*2,y,x+size,y+size)
			Vollkreis   (x+size,y+size-size*2/100+size/2,size/10)
			Volldreieck (x-size/10-1,y,x + size*2,y,x+size-size/10-1,y+size+size/2)
			Volldreieck (x,y,x+size*2+size/10+1,y,x+size+size/10+1,y+size+size/2)
			Volldreieck (x,y+size/10-1,x + size*2,y+size/10-1,x+size,y+size+size/2)
			
			
			Stiftfarbe(wr,wg,wb)
			Volldreieck (x,y,x + size*2,y,x+size,y-2*size)
			Vollkreis (x,y,size/10)
			Vollkreis (x + size*2,y,size/10)
			Vollkreis (x+size,y-2*size, size/10)
			Volldreieck (x-size/10-1,y,x + size*2,y,x+size-size/10-1,y-2*size)
			Volldreieck (x,y,x + size*2+size/10+1,y,x+size+size/10+1,y-2*size)
			Volldreieck (x,y+size/10-1,x + size*2,y+size/10-1,x+size,y-2*size-size/10)
			
			
			Stiftfarbe(pr,pg,pb)
			switch (*w).wert {
				case 1:
					Vollrechteck (x+size-size/20,y-size*11/10,size/10,size*4/5)
				case 2.:
					Vollrechteck (x+size*4/5,y-size*2/5,size*4/10,size/10)
					Vollrechteck (x+size*4/5,y-size*11/10,size*2/5,size/10)
					Vollrechteck (x+size*4/5,y-size*37/50,size*2/5,size/10)
					Vollrechteck (x+size*6/5,y-size,size/10,size*13/50)
					Vollrechteck (x+size*7/10,y-size*32/50,size/10,size*13/50)
				case 3:
					Vollrechteck (x+size*4/5,y-size*11/10,size*2/5,size/10)
					Vollrechteck (x+size*4/5,y-size*37/50,size*2/5,size/10)
					Vollrechteck (x+size*4/5,y-size*2/5,size*2/5,size/10)
					Vollrechteck (x+size*6/5,y-size,size/10,size*13/50)
					Vollrechteck (x+size*6/5,y-size*32/50,size/10,size*12/50)
				case 4:
					Vollrechteck (x+size*4/5,y-size*11/10,size/10,size*3/5)
					Vollrechteck (x+size*4/5,y-size*3/5,size*27/50,size/10)
					Vollrechteck (x+size*11/10,y-size*4/5,size/10,size*23/50)
				case 5:
					Vollrechteck (x+size*4/5,y-size*2/5,size*21/50,size/10)
					Vollrechteck (x+size*39/50,y-size*11/10,size*21/50,size/10)
					Vollrechteck (x+size*4/5,y-size*37/50,size*2/5,size/10)
					Vollrechteck (x+size*6/5,y-size*32/50,size/10,size*13/50)
					Vollrechteck (x+size*7/10,y-size,size/10,size*13/50)
				case 6:
					Vollrechteck (x+size*7/10,y-size,size/10,size*6/10)
					Vollrechteck (x+size*4/5,y-size*4/10,size*2/5,size/10)
					Vollrechteck (x+size*6/5,y-size*7/10,size/10,size*3/10)
					Vollrechteck (x+size*4/5,y-size*4/5,size*2/5,size/10)
					Vollrechteck (x+size*4/5,y-size*11/10,size*2/5,size/10)
					Vollrechteck (x+size*6/5,y-size,size/10,size/10)
				case 7:
					Vollrechteck (x+size*7/10,y-size*11/10,size*3/5,size/10)
					Vollrechteck (x+size*6/5,y-size,size/10,size/10)
					Vollrechteck (x+size*57/50,y-size*9/10,size/10,size/10)
					Vollrechteck (x+size*54/50,y-size*4/5,size/10,size/10)
					Vollrechteck (x+size*51/50,y-size*7/10,size/10,size/10)
					Vollrechteck (x+size*48/50,y-size*3/5,size/10,size/10)
					Vollrechteck (x+size*9/10,y-size*5/10,size/10,size/10)
					Vollrechteck (x+size*42/50,y-size*2/5,size/10,size/10)
				case 8:
					Vollrechteck (x+size*7/10,y-size,size/10,size*2/10)
					Vollrechteck (x+size*7/10,y-size*7/10,size/10,size*3/10)
					Vollrechteck (x+size*4/5,y-size*4/10,size*2/5,size/10)
					Vollrechteck (x+size*6/5,y-size*7/10,size/10,size*3/10)
					Vollrechteck (x+size*4/5,y-size*8/10,size*2/5,size/10)
					Vollrechteck (x+size*6/5,y-size,size/10,size*2/10)
					Vollrechteck (x+size*4/5,y-size*11/10,size*2/5,size/10)
				}
		
			case 12:	//Dodecaeder
				if (*w).highlight {
						Stiftfarbe(hr,hg,hb)
						Vollkreis(x,y,size*2/3+5)
				}
				Stiftfarbe(130,130,130)
				Vollkreis(x,y,size*2/3)
				
				Stiftfarbe(pr,pg,pb)
				Kreis(x,y,size*2/3)
				Linie(x,y-size/2,x,y-size*2/3)
				Linie(x+size*95/200,y-size*31/200,x+size*127/200,y-size*41/200)
				Linie(x-size*95/200,y-size*31/200,x-size*127/200,y-size*41/200)
				Linie(x+size*59/200,y+size*81/200,x+size*79/200,y+size*108/200)
				Linie(x-size*59/200,y+size*81/200,x-size*79/200,y+size*108/200)
				
				
				Stiftfarbe(wr,wg,wb)
				Volldreieck (x,y,x,y-size/2,x+size*95/200,y-size*31/200)
				Volldreieck (x,y,x,y-size/2,x-size*95/200,y-size*31/200)
				Volldreieck (x,y,x+size*59/200,y+size*81/200,x-size*59/200,y+size*81/200)
				Volldreieck (x,y,x-size*95/200,y-size*31/200,x-size*59/200,y+size*81/200)
				Volldreieck (x,y,x+size*95/200,y-size*31/200,x+size*59/200,y+size*81/200)
				var x12,y12,size12 uint16
				x12 = x-size*2/3+size*16/100
				y12 = y+size*2/3-size*32/100
				size12 = size/2
				Stiftfarbe(pr,pg,pb)
				switch (*w).wert {
					case 1:
						Vollrechteck (x12+size12-size12/20,y12-size12*11/10,size12/10,size12*4/5)
					case 2.:
						Vollrechteck (x12+size12*4/5,y12-size12*2/5,size12*4/10,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*37/50,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12*13/50)
						Vollrechteck (x12+size12*7/10,y12-size12*32/50,size12/10,size12*13/50)
					case 3:
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*37/50,size12*2/5,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*2/5,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12*13/50)
						Vollrechteck (x12+size12*6/5,y12-size12*32/50,size12/10,size12*12/50)
					case 4:
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12/10,size12*3/5)
						Vollrechteck (x12+size12*4/5,y12-size12*3/5,size12*27/50,size12/10)
						Vollrechteck (x12+size12*11/10,y12-size12*4/5,size12/10,size12*23/50)
					case 5:
						Vollrechteck (x12+size12*4/5,y12-size12*2/5,size12*21/50,size12/10)
						Vollrechteck (x12+size12*39/50,y12-size12*11/10,size12*21/50,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*37/50,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12*32/50,size12/10,size12*13/50)
						Vollrechteck (x12+size12*7/10,y12-size12,size12/10,size12*13/50)
					case 6:
						Vollrechteck (x12+size12*7/10,y12-size12,size12/10,size12*6/10)
						Vollrechteck (x12+size12*4/5,y12-size12*4/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12*7/10,size12/10,size12*3/10)
						Vollrechteck (x12+size12*4/5,y12-size12*4/5,size12*2/5,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12/10)
					case 7:
						Vollrechteck (x12+size12*7/10,y12-size12*11/10,size12*3/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12/10)
						Vollrechteck (x12+size12*57/50,y12-size12*9/10,size12/10,size12/10)
						Vollrechteck (x12+size12*54/50,y12-size12*4/5,size12/10,size12/10)
						Vollrechteck (x12+size12*51/50,y12-size12*7/10,size12/10,size12/10)
						Vollrechteck (x12+size12*48/50,y12-size12*3/5,size12/10,size12/10)
						Vollrechteck (x12+size12*9/10,y12-size12*5/10,size12/10,size12/10)
						Vollrechteck (x12+size12*42/50,y12-size12*2/5,size12/10,size12/10)
					case 8:
						Vollrechteck (x12+size12*7/10,y12-size12,size12/10,size12*2/10)
						Vollrechteck (x12+size12*7/10,y12-size12*7/10,size12/10,size12*3/10)
						Vollrechteck (x12+size12*4/5,y12-size12*4/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12*7/10,size12/10,size12*3/10)
						Vollrechteck (x12+size12*4/5,y12-size12*8/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12*2/10)
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12*2/5,size12/10)
					case 9:
						Vollrechteck (x12+size12*7/10,y12-size12,size12/10,size12*2/10)
						Vollrechteck (x12+size12*4/5,y12-size12*4/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*8/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12*6/10)
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12*2/5,size12/10)
					case 10:
						//1
						x12 = x12-size*15/100
						Vollrechteck (x12+size12-size12/20,y12-size12*11/10,size12/10,size12*4/5)
						//0
						x12 = x12+size*28/100
						Vollrechteck (x12+size12*7/10,y12-size12,size12/10,size12*6/10)
						Vollrechteck (x12+size12*4/5,y12-size12*4/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12*6/10)
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12*2/5,size12/10)
					case 11:
						//1
						x12 = x12-size*15/100
						Vollrechteck (x12+size12-size12/20,y12-size12*11/10,size12/10,size12*4/5)
						//1
						x12 = x12+size*28/100
						Vollrechteck (x12+size12-size12/20,y12-size12*11/10,size12/10,size12*4/5)
					case 12:
						//1
						x12 = x12-size*15/100
						Vollrechteck (x12+size12-size12/20,y12-size12*11/10,size12/10,size12*4/5)
						//2
						x12 = x12+size*28/100
						Vollrechteck (x12+size12*4/5,y12-size12*2/5,size12*4/10,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*37/50,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12*13/50)
						Vollrechteck (x12+size12*7/10,y12-size12*32/50,size12/10,size12*13/50)
					}
			
			case 20:	//Icosanoeder
				if (*w).highlight {
						Stiftfarbe(hr,hg,hb)
						Volldreieck(x,y+size*24/10,x+size*43/10,y,x+size*43/10,y-7)											//Dreieck 1
						Volldreieck(x,y+size*24/10,x+size*43/10,y-7,x-5,y+size*24/10-5)										//Dreieck 2
						Volldreieck(x,y+size*24/10,x,y+size*73/10,x-5,y+size*24/10-5)										//Dreisck 3
						Volldreieck(x,y+size*73/10,x-5,y+size*24/10-5,x-5,y+size*73/10+5)									//Dreieck 4
						Volldreieck(x,y+size*73/10,x+size*43/10,y+size*98/10,x+size*43/10,y+size*98/10+7)					//Dreieck 5
						Volldreieck(x,y+size*73/10,x-5,y+size*73/10+5,x+size*43/10,y+size*98/10+7)							//Dreieck 6
						Volldreieck(x+size*43/10,y+size*98/10,x+size*43/10,y+size*98/10+7,x+size*86/10+5,y+size*73/10+5)	//Dreieck 7
						Volldreieck(x+size*43/10,y+size*98/10,x+size*86/10,y+size*73/10,x+size*86/10+5,y+size*73/10+5)		//Dreieck 8
						Volldreieck(x+size*86/10,y+size*73/10,x+size*86/10,y+size*24/10,x+size*86/10+5,y+size*24/10-5)		//Dreieck 9
						Volldreieck(x+size*86/10,y+size*73/10,x+size*86/10+5,y+size*73/10+5,x+size*86/10+5,y+size*24/10-5)	//Dreieck 10
						Volldreieck(x+size*86/10,y+size*24/10,x+size*86/10+5,y+size*24/10-5,x+size*43/10,y-7)				//Dreieck 11
						Volldreieck(x+size*43/10,y,x+size*86/10,y+size*24/10,x+size*43/10,y-7)								//Dreieck 12
						
				}
				Stiftfarbe(130,130,130)
				Volldreieck(x,y+size*24/10,x+size*43/10,y,x+size*43/10,y+size*24/10)						//Dreieck 1
				Volldreieck(x+size*43/10,y,x+size*43/10,y+size*24/10,x+size*86/10,y+size*24/10)				//Dreieck 2
				Volldreieck(x,y+size*24/10,x+size*43/10,y+size*24/10,x+size*21/10,y+size*61/10)				//Dreisck 3
				//Volldreieck(x+size*43/10,y+size*24/10,x+size*21/10,y+size*61/10,x+size*65/10,y+size*61/10)	//Dreieck 4
				Volldreieck(x+size*43/10,y+size*24/10,x+size*86/10,y+size*24/10,x+size*65/10,y+size*61/10)	//Dreieck 5
				Volldreieck(x,y+size*24/10,x+size*21/10,y+size*61/10,x,y+size*73/10)						//Dreieck 6
				Volldreieck(x+size*86/10,y+size*24/10,x+size*86/10,y+size*73/10,x+size*65/10,y+size*61/10)	//Dreieck 7
				Volldreieck(x,y+size*73/10,x+size*21/10,y+size*61/10,x+size*43/10,y+size*98/10)				//Dreieck 8
				Volldreieck(x+size*21/10,y+size*61/10,x+size*65/10,y+size*61/10,x+size*43/10,y+size*98/10)	//Dreieck 9
				Volldreieck(x+size*65/10,y+size*61/10,x+size*86/10,y+size*73/10,x+size*43/10,y+size*98/10)	//Dreieck 10
				
				Stiftfarbe(wr,wg,wb)
				Volldreieck(x+size*43/10,y+size*24/10,x+size*21/10,y+size*61/10,x+size*65/10,y+size*61/10)	//Dreieck 4
				
				Stiftfarbe(pr,pg,pb)
				Dreieck(x,y+size*24/10,x+size*43/10,y,x+size*43/10,y+size*24/10)						//Dreieck 1
				Dreieck(x+size*43/10,y,x+size*43/10,y+size*24/10,x+size*86/10,y+size*24/10)				//Dreieck 2
				Dreieck(x,y+size*24/10,x+size*43/10,y+size*24/10,x+size*21/10,y+size*61/10)				//Dreisck 3
				Dreieck(x+size*43/10,y+size*24/10,x+size*21/10,y+size*61/10,x+size*65/10,y+size*61/10)	//Dreieck 4
				Dreieck(x+size*43/10,y+size*24/10,x+size*86/10,y+size*24/10,x+size*65/10,y+size*61/10)	//Dreieck 5
				Dreieck(x,y+size*24/10,x+size*21/10,y+size*61/10,x,y+size*73/10)						//Dreieck 6
				Dreieck(x+size*86/10,y+size*24/10,x+size*86/10,y+size*73/10,x+size*65/10,y+size*61/10)	//Dreieck 7
				Dreieck(x,y+size*73/10,x+size*21/10,y+size*61/10,x+size*43/10,y+size*98/10)				//Dreieck 8
				Dreieck(x+size*21/10,y+size*61/10,x+size*65/10,y+size*61/10,x+size*43/10,y+size*98/10)	//Dreieck 9
				Dreieck(x+size*65/10,y+size*61/10,x+size*86/10,y+size*73/10,x+size*43/10,y+size*98/10)	//Dreieck 10
				
				
				
				var x12,y12,size12 uint16
				x12 = x+size*233/100
				y12 = y+size*65/10
				size12 = size*2
				Stiftfarbe(pr,pg,pb)
				switch (*w).wert {
					case 1:
						Vollrechteck (x12+size12-size12/20,y12-size12*11/10,size12/10,size12*4/5)
					case 2.:
						Vollrechteck (x12+size12*4/5,y12-size12*2/5,size12*4/10,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*37/50,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12*13/50)
						Vollrechteck (x12+size12*7/10,y12-size12*32/50,size12/10,size12*13/50)
					case 3:
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*37/50,size12*2/5,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*2/5,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12*13/50)
						Vollrechteck (x12+size12*6/5,y12-size12*32/50,size12/10,size12*12/50)
					case 4:
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12/10,size12*3/5)
						Vollrechteck (x12+size12*4/5,y12-size12*3/5,size12*27/50,size12/10)
						Vollrechteck (x12+size12*11/10,y12-size12*4/5,size12/10,size12*23/50)
					case 5:
						Vollrechteck (x12+size12*4/5,y12-size12*2/5,size12*21/50,size12/10)
						Vollrechteck (x12+size12*39/50,y12-size12*11/10,size12*21/50,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*37/50,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12*32/50,size12/10,size12*13/50)
						Vollrechteck (x12+size12*7/10,y12-size12,size12/10,size12*13/50)
					case 6:
						Vollrechteck (x12+size12*7/10,y12-size12,size12/10,size12*6/10)
						Vollrechteck (x12+size12*4/5,y12-size12*4/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12*7/10,size12/10,size12*3/10)
						Vollrechteck (x12+size12*4/5,y12-size12*4/5,size12*2/5,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12/10)
					case 7:
						Vollrechteck (x12+size12*7/10,y12-size12*11/10,size12*3/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12/10)
						Vollrechteck (x12+size12*57/50,y12-size12*9/10,size12/10,size12/10)
						Vollrechteck (x12+size12*54/50,y12-size12*4/5,size12/10,size12/10)
						Vollrechteck (x12+size12*51/50,y12-size12*7/10,size12/10,size12/10)
						Vollrechteck (x12+size12*48/50,y12-size12*3/5,size12/10,size12/10)
						Vollrechteck (x12+size12*9/10,y12-size12*5/10,size12/10,size12/10)
						Vollrechteck (x12+size12*42/50,y12-size12*2/5,size12/10,size12/10)
					case 8:
						Vollrechteck (x12+size12*7/10,y12-size12,size12/10,size12*2/10)
						Vollrechteck (x12+size12*7/10,y12-size12*7/10,size12/10,size12*3/10)
						Vollrechteck (x12+size12*4/5,y12-size12*4/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12*7/10,size12/10,size12*3/10)
						Vollrechteck (x12+size12*4/5,y12-size12*8/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12*2/10)
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12*2/5,size12/10)
					case 9:
						Vollrechteck (x12+size12*7/10,y12-size12,size12/10,size12*2/10)
						Vollrechteck (x12+size12*4/5,y12-size12*4/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*8/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12*6/10)
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12*2/5,size12/10)
					case 10:
						//1
						x12 = x12-size*60/100
						Vollrechteck (x12+size12-size12/20,y12-size12*11/10,size12/10,size12*4/5)
						//0
						x12 = x12+size*112/100
						Vollrechteck (x12+size12*7/10,y12-size12,size12/10,size12*6/10)
						Vollrechteck (x12+size12*4/5,y12-size12*4/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12*6/10)
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12*2/5,size12/10)
					case 11:
						//1
						x12 = x12-size*60/100
						Vollrechteck (x12+size12-size12/20,y12-size12*11/10,size12/10,size12*4/5)
						//1
						x12 = x12+size*112/100
						Vollrechteck (x12+size12-size12/20,y12-size12*11/10,size12/10,size12*4/5)
					case 12:
						//1
						x12 = x12-size*60/100
						Vollrechteck (x12+size12-size12/20,y12-size12*11/10,size12/10,size12*4/5)
						//2
						x12 = x12+size*112/100
						Vollrechteck (x12+size12*4/5,y12-size12*2/5,size12*4/10,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*37/50,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12*13/50)
						Vollrechteck (x12+size12*7/10,y12-size12*32/50,size12/10,size12*13/50)
					case 13:
						//1
						x12 = x12-size*60/100
						Vollrechteck (x12+size12-size12/20,y12-size12*11/10,size12/10,size12*4/5)
						//3
						x12 = x12+size*112/100
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*37/50,size12*2/5,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*2/5,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12*13/50)
						Vollrechteck (x12+size12*6/5,y12-size12*32/50,size12/10,size12*12/50)
					case 14:
						//1
						x12 = x12-size*60/100
						Vollrechteck (x12+size12-size12/20,y12-size12*11/10,size12/10,size12*4/5)
						//4
						x12 = x12+size*112/100
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12/10,size12*3/5)
						Vollrechteck (x12+size12*4/5,y12-size12*3/5,size12*27/50,size12/10)
						Vollrechteck (x12+size12*11/10,y12-size12*4/5,size12/10,size12*23/50)
					case 15:
						//1
						x12 = x12-size*60/100
						Vollrechteck (x12+size12-size12/20,y12-size12*11/10,size12/10,size12*4/5)
						//5
						x12 = x12+size*112/100
						Vollrechteck (x12+size12*4/5,y12-size12*2/5,size12*21/50,size12/10)
						Vollrechteck (x12+size12*39/50,y12-size12*11/10,size12*21/50,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*37/50,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12*32/50,size12/10,size12*13/50)
						Vollrechteck (x12+size12*7/10,y12-size12,size12/10,size12*13/50)
					case 16:
						//1
						x12 = x12-size*60/100
						Vollrechteck (x12+size12-size12/20,y12-size12*11/10,size12/10,size12*4/5)
						//6
						x12 = x12+size*112/100
						Vollrechteck (x12+size12*7/10,y12-size12,size12/10,size12*6/10)
						Vollrechteck (x12+size12*4/5,y12-size12*4/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12*7/10,size12/10,size12*3/10)
						Vollrechteck (x12+size12*4/5,y12-size12*4/5,size12*2/5,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12/10)
					case 17:
						//1
						x12 = x12-size*60/100
						Vollrechteck (x12+size12-size12/20,y12-size12*11/10,size12/10,size12*4/5)
						//7
						x12 = x12+size*112/100
						Vollrechteck (x12+size12*7/10,y12-size12*11/10,size12*3/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12/10)
						Vollrechteck (x12+size12*57/50,y12-size12*9/10,size12/10,size12/10)
						Vollrechteck (x12+size12*54/50,y12-size12*4/5,size12/10,size12/10)
						Vollrechteck (x12+size12*51/50,y12-size12*7/10,size12/10,size12/10)
						Vollrechteck (x12+size12*48/50,y12-size12*3/5,size12/10,size12/10)
						Vollrechteck (x12+size12*9/10,y12-size12*5/10,size12/10,size12/10)
						Vollrechteck (x12+size12*42/50,y12-size12*2/5,size12/10,size12/10)
					case 18:
						//1
						x12 = x12-size*60/100
						Vollrechteck (x12+size12-size12/20,y12-size12*11/10,size12/10,size12*4/5)
						//8
						x12 = x12+size*112/100
						Vollrechteck (x12+size12*7/10,y12-size12,size12/10,size12*2/10)
						Vollrechteck (x12+size12*7/10,y12-size12*7/10,size12/10,size12*3/10)
						Vollrechteck (x12+size12*4/5,y12-size12*4/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12*7/10,size12/10,size12*3/10)
						Vollrechteck (x12+size12*4/5,y12-size12*8/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12*2/10)
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12*2/5,size12/10)
					case 19:
						//1
						x12 = x12-size*60/100
						Vollrechteck (x12+size12-size12/20,y12-size12*11/10,size12/10,size12*4/5)
						//9
						x12 = x12+size*112/100
						Vollrechteck (x12+size12*7/10,y12-size12,size12/10,size12*2/10)
						Vollrechteck (x12+size12*4/5,y12-size12*4/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*8/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12*6/10)
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12*2/5,size12/10)
					case 20:
						//2
						x12 = x12-size*68/100
						Vollrechteck (x12+size12*4/5,y12-size12*2/5,size12*4/10,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*4/5,y12-size12*37/50,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12*13/50)
						Vollrechteck (x12+size12*7/10,y12-size12*32/50,size12/10,size12*13/50)
						//0
						x12 = x12+size*128/100
						Vollrechteck (x12+size12*7/10,y12-size12,size12/10,size12*6/10)
						Vollrechteck (x12+size12*4/5,y12-size12*4/10,size12*2/5,size12/10)
						Vollrechteck (x12+size12*6/5,y12-size12,size12/10,size12*6/10)
						Vollrechteck (x12+size12*4/5,y12-size12*11/10,size12*2/5,size12/10)
					
					}
			default:
				panic("Sie haben ein Würfel gewählt, der nicht von dem Würfelpaket gezeichnet werden kann!! Wählen Sie New(2,0), New(4,0), New(6,0), New(8,0), New(12,0) oder New(20,0) oder ihre Invertierung!!!") 
	}
}


		

