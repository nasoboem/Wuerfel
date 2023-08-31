/* Autor: St. Schmidt
 * Datum: 27.03.2017
 * Zweck: Zufallszahlen */

package zufallszahlen

import ( "math/rand" ; "time" )

// Vor.: keine
// Eff.: Der Zufallszahlengenerator ist mit keim initialisiert.
func Initialisieren (keim int64) {
	rand.Seed (keim)
}

// Vor.: keine
// Eff.: Der Zufallszahlengenerator ist mit einem aus der Systemzeit
//       abgeleiteten Wert neu initialisiert.
func Randomisieren () {
	rand.Seed(time.Now().UnixNano())
}

// Vor.: -2^63 <= a <= b <= 2^63-1
// Erg.: Eine ganze Zufallszahl aus dem Intervall [a,b] ist
//       geliefert.
func Zufallszahl (a,b int64) int64 {
	return a + rand.Int63n(b-a+1)
}
