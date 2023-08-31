package main

import (
		"würfel/gfx"
		"würfel/wuerfel"
		)


func main () {
	var w1,w2,w3,w4,w5,w6,p1,p2,p3,p4,p5,p6 wuerfel.Wuerfel
	w1 = wuerfel.New(2,0)
	w2 = wuerfel.New(4,0)
	w3 = wuerfel.New(6,0)
	w4 = wuerfel.New(8,0)
	w5 = wuerfel.New(12,0)
	w6 = wuerfel.New(20,0)
	p1 = wuerfel.New(0,2)
	p2 = wuerfel.New(0,4)
	p3 = wuerfel.New(0,6)
	p4 = wuerfel.New(0,8)
	p5 = wuerfel.New(0,12)
	p6 = wuerfel.New(0,20)
	var liste []wuerfel.Wuerfel
	liste = append(liste,w1,w2,w3,w4,w5,w6,p1,p2,p3,p4,p5,p6)
	
	//MausLesen1 () (taste uint8, status int8, mausX, mausY uint16)
	var taste uint8
	var status int8
	var mx,my uint16
	gfx.Fenster(1600,1000)
	
	gfx.UpdateAus()
	gfx.Stiftfarbe(10,255,5)
	gfx.Cls()
	for i:=0;i<len(liste);i++ {
		liste[i].Draw(150+uint16(i%6)*200,200+uint16(i/6)*400,40)
	}
	gfx.UpdateAn()
	
	for {
		taste,status,mx,my = gfx.MausLesen1()
		for i:=0;i<len(liste);i++ {
			if liste[i].PunktgehörtzumWuerfel(mx,my) {
				liste[i].SetzeHighlight(true)
			} else {
				liste[i].SetzeHighlight(false)
			}
		}
		for i:=0;i<len(liste);i++ {
			if liste[i].PunktgehörtzumWuerfel(mx,my) && taste == 1 && status == 1 {
				if liste[i].GibSeiten()==0 {
					if i<6 {
						liste[i],liste[i+6] = liste[i+6],liste[i]
						break
					}else{
						liste[i],liste[i-6] = liste[i-6],liste[i]
						break
					}
				}else{
					liste[i].Wuerfeln()
					break
				}
			}
		}
		gfx.UpdateAus()
		gfx.Stiftfarbe(10,255,5)
		gfx.Cls()
		for i:=0;i<len(liste);i++ {
			liste[i].Draw(150+uint16(i%6)*200,200+uint16(i/6)*400,40)
		}
		gfx.UpdateAn()
	}
}

