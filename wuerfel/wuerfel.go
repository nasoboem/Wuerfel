package wuerfel

type Wuerfel interface {
	
	//Vor.: -
	//Erg.:Ein Würfel mit der gegeben Seitenzahl ist zurückgegeben. Bisher implementiert: 0, 2, 4, 6, 8, 12
	//Für den Schatten eines Würfels geben Sie die Seitenzahl 0 an und dann die Seitenzahl des Würfels für den der Schatten sein soll.
	//Eff.: Würfelwert ist 0 (entspricht ungewürfelt).
	//New (seiten,schatten uint) Wuerfel  -- Bsp.: 6-seitiger Würfel: New(6,0); schatten eines 6-Seitigen Würfels: New(0,6)
	
	//Vor.: -
	//Erg.: Der aktuelle Wert des Würfels ist geliefert.
	GibWert () uint
	
	//Vor.: -
	//Erg.: Die Anzahl der Seiten des Würfels (Max. Würfelwert) ist geliefert.
	GibSeiten () uint
	
	//Vor.: -
	//Erg.: Die Farbe des Würfels in r,g,b ist zurückgegeben. (Default: 255,255,255 - weiß)
	GibWuerfelFarbe () (r,g,b uint8)
	
	//Vor.: -
	//Erg.: Die Farbe der Punkte in r,g,b ist zurückgegeben. (Default: 0,0,0 - schwarz)
	GibPunktFarbe () (r,g,b uint8)
	
	//Vor.: Der Würfel wurde gezeichnet. (Draw-Funktion) (Default: x,y = 0,0) 
	//Erg.: Die Position, an der der Würfel gezeichnet wurde ist zurückgegeben.
	GibPosition () (x,y uint16)
	
	//Vor.: Der Würfel wurde gezeichnet. (Draw-Funktion) (Default: size = 0)
	//Erg.: Die Größe mit der der Würfel gezeichnet wurde ist zurückgegeben. (Draw-Funktion) (Default: size = 0)
	GibGroesse () (size uint16)
	
	//Vor.: -
	// Erg.: Die Highlight-Farbe ist in r,g,b zurückgegeben. (Default: 255,0,0 - rot)
	GibHighlightFarbe() (r,g,b uint8)
	
	//Vor.: Für den vollen Funktionsumfang wählen Sie einen der implementierten Würfel: 0(Schatten/Platzhalter),2,4,6,8,12, sonst seiten <=0.
	//Eff.: Der Würfel besitzt die übergeben Seitenzahl und den Würfelwert = 0
	SetzeSeiten (seiten uint)
	
	//Vor.: -
	//Eff.:	Würfelwert besitzt den gegeben Wert, wenn Wert <= Seitenzahl des Würfels ist,
	//		andernfalls passiert nichts.
	SetzeWert (wert uint)
	
	//Vor.: -
	//Eff.: Würfelwert besitzt den gegeben Wert, wenn Wert <= Seitenzahl des Würfels ist,
	//		wird der Wert des Würfels maximal (Seitenzahl).
	SetzeWertb (wert uint)
	
	//Vor.: Wert <= Seitenzahl des Würfels
	//Eff.: Würfelwert besitzt den gegeben Wert.
	SetzeWertc (wert uint)
	
	//Vor.: -
	//Eff.: Der Würfel besitzt die mit r,g,b gegebene Farbe.
	SetzeWuerfelFarbe (r,g,b uint8)
	
	//Vor.: -
	//Eff.: Die Punkte besitzten die mit r,g,b gegebene Farbe.
	SetzePunktFarbe (r,g,b uint8)
	
	//Vor.: -
	//Eff.: Die Highlight farbe hat jetzt die gegeben Werte r,g,b.
	SetzeHighlightFarbe (r,g,b uint8)
	
	//Vor.: -
	//Eff.: Der Würfel hat den Wert 0 - ungewürfelt, es werden keine Punkte angezeigt.
	Zuruecksetzen ()
	
	//Vor.: -
	//Eff.: Bei true wird das Highlight in der Draw-Funktion angezeigt und bei false nicht.	
	SetzeHighlight (highlight bool)
	
	//Vor.: -
	//Erg.: Der Würfelwert wird als String zurückgegeben.
	String () string
	
	//Vor.: -
	//Eff.: Der Würfel erhält einen neuen zufälligen Würfelwert, zwischen 1 und der Seitenanzahl (max. Würfelwert).	
	Wuerfeln ()
	
	//Vor.: Würfel muss vorher gezeichnet sein
	//Erg.: True ist geliefert, wenn der Punkt zum Würfel gehört.
	PunktgehörtzumWuerfel (xp,yp uint16) bool
	
	//Vor.: Ein gfx-Fenster ist geöffnet.
	//Eff.:	1) Die Werte x,y und size werden im Würfel gespeichert.
	//		2) Der Würfel ist den Werten x,y und size entsprechend gezeichnet.
	Draw(x,y,size uint16)
}

