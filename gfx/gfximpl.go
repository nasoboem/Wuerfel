// Autor: (c) Stefan Schmidt (Kontakt: St.Schmidt@online.de);
// Datum: 07.03.2016 ; letzte Änderung: 26.04.2021;
// Zweck: Grafik- und Soundausgabe und Eingabe per Tastatur und Maus
// mit Go unter Windows und unter Linux
package gfx
/* Letzte Änderungen:
- 26.04.2021  Veränderung bei 'SchreibeFont': Leerstring mit Sonderbehandlung, da sonst
              C.Write ein 'Sementation Fault' auslöst
- 06.11.2020  Bug bei 'Clipboard_einfuegen' bzgl. der Transparenz entfernt
- 18.10.2020  alle "docstrings" in die Impl. kopiert, so dass nun 
              'go doc gfx.<FktName>' die Spezifikation liefert, 
              Bug in 'FensterAus' entfernt,
              neue Funktion 'Fenstertitel', mit der in der Titelzeile des Fensters
              ein eigener Fenstertitel festgelegt werden kann,
              neue Funktionen 'LadeBildMitColorKey' und 'Clipboard_einfuegenMitColorKey',
              bei der Pixel einer bestimmten
              Farbe transparent dargestellt werden (gut für "Sprites"!),
              neue Funktion 'Transparenz', damit man sich überdeckende Grafik-
              objekte erkennen kann
- 01.10.2020  'Bug' entfernt: Mit 'defer' angemeldete Funktionsaufrufe
              wurden nach dem Schließen des Fensters mit Klick auf das x links oben
              nicht mehr ausgeführt. 
- 01.09.2019  Spezifikationsfehler korrigiert
- 08.05.2019  neue Funktionen, um (Klavier-) Noten spielen zu können
              inkl. Hüllkurven und klanganpassungen 
- 07.03.2019  Rechtschreibkorrekturen in der Spezifikation 
- 03.03.2019: Die Funktion 'SetzeFont' liefert nun einen Rückgabewert,
              der den Erfolg/Misserfolg angibt.
- 07.10.2017: neue Funktion 'Tastaturzeichen' für deutsche Tastaturbelegung!
- 07.10.2017: 'Bug' in Funktion 'Cls()' entfernt - KEIN FLACKERN MEHR 
              bei 'double-buffering' mit UpdateAus() und UpdateAn()
*/

/*
#cgo LDFLAGS: -lSDL -lSDL_gfx -lSDL_ttf
#include <SDL/SDL.h>
#include <SDL/SDL_ttf.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include <SDL/SDL_gfxPrimitives.h>

// Structure for loaded sounds. 
typedef struct sound_s {
    Uint8 *samples;		// raw PCM sample data 
    Uint32 length;		// size of sound data in bytes 
} sound_t, *sound_p;

// Structure for a currently playing sound. 
typedef struct playing_s {
    int active;                 // 1 if this sound should be played 
    sound_p sound;              // sound data to play 
    Uint32 position;            // current position in the sound buffer 
} playing_t, *playing_p;

// Array for all active sound effects. 
#define MAX_PLAYING_SOUNDS      10 
playing_t playing[MAX_PLAYING_SOUNDS];

// The higher this is, the louder each currently playing sound will be.
// However, high values may cause distortion if too many sounds are
// playing. Experiment with this. 
#define VOLUME_PER_SOUND        SDL_MIX_MAXVOLUME / 2

static SDL_Surface *screen;
static SDL_Surface *archiv;
static SDL_Surface *clipboard = NULL;
static Uint8 updateOn = 1;
static Uint8 red,green,blue,alpha;
static Uint8 ck_red, ck_green, ck_blue;
static Uint8 colorkey;
static SDL_Event event; 
static Uint8 gedrueckt;
static Uint16 taste,tiefe;
static Uint8 tasteLesen = 0;
static Uint8 tastaturpuffer = 0;
static Uint32 t_puffer[256];
static Uint8 t_pufferkopf;
static Uint8 t_pufferende;
static Uint16 mausX, mausY;
static Uint8 mausLesen = 0;
static Uint8 mausTaste;
static Uint8 mauspuffer = 0;
static Uint32 m_puffer[256];
static Uint8 m_pufferkopf;
static Uint8 m_pufferende;
static Uint8 fensteroffen = 0;
static Uint8 fensterzu = 1;
static char aktFont[256];
static int aktFontSize;
static SDL_AudioSpec desired, obtained; // Audio format specifications.
static sound_t s[10];                   // Our loaded sounds and their formats. 
   
//------------------------------------------------------------------
// This function is called by SDL whenever the sound card
// needs more samples to play. It might be called from a
// separate thread, so we should be careful what we touch. 
static void AudioCallback(void *user_data, Uint8 *audio, int length)
{
    int i;
    // Avoid compiler warning. 
    user_data += 0;
    // Clear the audio buffer so we can mix samples into it. 
    memset(audio, 0, length);
    // Mix in each sound. 
    for (i = 0; i < MAX_PLAYING_SOUNDS; i++) {
	  if (playing[i].active) {
	    Uint8 *sound_buf;
	    Uint32 sound_len;
	    // Locate this sound's current buffer position. 
	    sound_buf = playing[i].sound->samples;
	    sound_buf += playing[i].position;
	    // Determine the number of samples to mix. 
	    if ((playing[i].position + length) > playing[i].sound->length) {
		sound_len = playing[i].sound->length - playing[i].position;
	    } else {
		sound_len = length;
	    }
	    // Mix this sound into the stream. 
	    SDL_MixAudio(audio, sound_buf, sound_len, VOLUME_PER_SOUND);
	    // Update the sound buffer's position. 
	    playing[i].position += length;
	    // Have we reached the end of the sound? 
	    if (playing[i].position >= playing[i].sound->length) {
	    free(s[i].samples);      //zugehörigen Soundstruktur-Samplespeicher wieder freigeben 
		playing[i].active = 0;	 // und anschließend als inaktiv markieren
	    }
	  }
    }
}
//----------------------------------------------------------------
// This function loads a sound with SDL_LoadWAV and converts
// it to the specified sample format. Returns 0 on success
// and 1 on failure. 
static int LoadAndConvertSound(char *filename, SDL_AudioSpec *spec,
			sound_p sound)
{
    SDL_AudioCVT cvt;           // audio format conversion structure 
    SDL_AudioSpec loaded;       // format of the loaded data 
    Uint8 *new_buf;
    // Load the WAV file in its original sample format. 
    if (SDL_LoadWAV(filename,
		    &loaded, &sound->samples, &sound->length) == NULL) {
	printf("Unable to load sound: %s\n", SDL_GetError());
	return 1;
    }
    // Build a conversion structure for converting the samples.
    // This structure contains the data SDL needs to quickly
    // convert between sample formats. 
    if (SDL_BuildAudioCVT(&cvt, loaded.format,
			  loaded.channels,
			  loaded.freq,
			  spec->format, spec->channels, spec->freq) < 0) {
	// printf("Unable to convert sound: %s\n", SDL_GetError());
	return 1;
    }
    // Since converting PCM samples can result in more data
    //   (for instance, converting 8-bit mono to 16-bit stereo),
    //   we need to allocate a new buffer for the converted data.
    //   Fortunately SDL_BuildAudioCVT supplied the necessary
    //   information. 
    cvt.len = sound->length;
    new_buf = (Uint8 *) malloc(cvt.len * cvt.len_mult);
    if (new_buf == NULL) {
	//printf("Memory allocation failed.\n");
	SDL_FreeWAV(sound->samples);
	return 1;
    }
    // Copy the sound samples into the new buffer.
    memcpy(new_buf, sound->samples, sound->length);
    // Perform the conversion on the new buffer. 
    cvt.buf = new_buf;
    if (SDL_ConvertAudio(&cvt) < 0) {
	//printf("Audio conversion error: %s\n", SDL_GetError());
	free(new_buf);
	SDL_FreeWAV(sound->samples);
	return 1;
    }
    // Swap the converted data for the original. 
    SDL_FreeWAV(sound->samples);
    sound->samples = new_buf;
    sound->length = sound->length * cvt.len_mult;
    // Success! 
    //printf("'%s' was loaded and converted successfully.\n", filename);
    return 0;
}
//----------------------------------------------------------------
// Diese Funktion übernimmt eine Bytefolge aus dem RAM ab der Adresse addr
// mit der Länge laenge, die dem Inhalt einer WAV-DAtei entspricht und konvertiert
// sie , damit es abgespielt werden kann. Die Funktion liefert 0 bei Erfolg
// 1 bei Misserfolg. 
static int LadeUndKonvertiereRAMWAV(const void* addr, int laenge, SDL_AudioSpec *spec,
			sound_p sound)
{
    SDL_AudioCVT cvt;           // audio format conversion structure 
    SDL_AudioSpec loaded;       // format of the loaded data 
    Uint8 *new_buf;
    // Lade 'RAMWAV' im Originalformat:  
    if (SDL_LoadWAV_RW(SDL_RWFromConstMem(addr,laenge),0,
		    &loaded, &sound->samples, &sound->length) == NULL) {
	printf("Unable to load sound: %s\n", SDL_GetError());
	return 1;
    }
    // Build a conversion structure for converting the samples.
    // This structure contains the data SDL needs to quickly
    // convert between sample formats. 
    if (SDL_BuildAudioCVT(&cvt, loaded.format,
			  loaded.channels,
			  loaded.freq,
			  spec->format, spec->channels, spec->freq) < 0) {
	// printf("Unable to convert sound: %s\n", SDL_GetError());
	return 1;
    }
    // Since converting PCM samples can result in more data
    //   (for instance, converting 8-bit mono to 16-bit stereo),
    //   we need to allocate a new buffer for the converted data.
    //   Fortunately SDL_BuildAudioCVT supplied the necessary
    //   information. 
    cvt.len = sound->length;
    new_buf = (Uint8 *) malloc(cvt.len * cvt.len_mult);
    if (new_buf == NULL) {
	//printf("Memory allocation failed.\n");
	SDL_FreeWAV(sound->samples);
	return 1;
    }
    // Copy the sound samples into the new buffer.
    memcpy(new_buf, sound->samples, sound->length);
    // Perform the conversion on the new buffer. 
    cvt.buf = new_buf;
    if (SDL_ConvertAudio(&cvt) < 0) {
	//printf("Audio conversion error: %s\n", SDL_GetError());
	free(new_buf);
	SDL_FreeWAV(sound->samples);
	return 1;
    }
    // Swap the converted data for the original. 
    SDL_FreeWAV(sound->samples);
    sound->samples = new_buf;
    sound->length = sound->length * cvt.len_mult;
    // Success! 
    //printf("'%s' was loaded and converted successfully.\n", filename);
    return 0;
}
//-----------------------------------------------------------------
static int LoadAndPlaySound (char *filename) 
{
	int i;
	//Finde einen freien Index (Bereich 0 <= index < MAX_PLAYING_SOUND
	for (i = 0; i < MAX_PLAYING_SOUNDS; i++) {
	if (playing[i].active == 0)
	    break;
    } 
    if (i == MAX_PLAYING_SOUNDS)
	return 1; //Fehler: Es werden schon die max. Anzahl an Dateien abgespielt.

	//Lade und konvertiere den Sound in die entsprechende Soundstruktur
	if (LoadAndConvertSound(filename, &obtained, &s[i]) != 0) {
	  return 2; //Laden fehlgeschlagen!
    }
    //Abspielen starten
    // The 'playing' structures are accessed by the audio callback,
    // so we should obtain a lock before we access them. 
    SDL_LockAudio();
    playing[i].active = 1;
    playing[i].sound = &s[i];
    playing[i].position = 0;
    SDL_UnlockAudio();
    return 0;
}    
//-----------------------------------------------------------------
static int LadeUndSpieleNote (const void* addr, int laenge) 
{
	int i;
	//Finde einen freien Index (Bereich 0 <= index < MAX_PLAYING_SOUND
	for (i = 0; i < MAX_PLAYING_SOUNDS; i++) {
	if (playing[i].active == 0)
	    break;
    } 
    if (i == MAX_PLAYING_SOUNDS)
	return 1; //Fehler: Es werden schon die max. Anzahl an Dateien abgespielt.

	//Lade und konvertiere den Sound in die entsprechende Soundstruktur
	if (LadeUndKonvertiereRAMWAV(addr, laenge, &obtained, &s[i]) != 0) {
	  return 2; //Laden fehlgeschlagen!
    }
    //Abspielen starten
    // The 'playing' structures are accessed by the audio callback,
    // so we should obtain a lock before we access them. 
    SDL_LockAudio();
    playing[i].active = 1;
    playing[i].sound = &s[i];
    playing[i].position = 0;
    SDL_UnlockAudio();
    return 0;
}    
//------------------------------------------------------------------
static int setFont (char *fontfile, int groesse) {
	strcpy (aktFont,fontfile);
	aktFontSize = groesse;
	TTF_Font *font = TTF_OpenFont(aktFont, aktFontSize);
	if (!font) {
	  //printf("TTF_OpenFont: %s\n", TTF_GetError());
	  return 1;
    }
    TTF_CloseFont(font);
	return 0;
}
//------------------------------------------------------------------
static char *getFont () {
	return aktFont;
}
//------------------------------------------------------------------
static int write (Sint16 x, Sint16 y, char *text) {
	TTF_Font *font = TTF_OpenFont(aktFont, aktFontSize);
	if (!font) {
	  //printf("TTF_OpenFont: %s\n", TTF_GetError());
	  return 1;
    }
	SDL_Color clrFg = {red,green,blue,alpha};  
	SDL_Surface *sText = TTF_RenderUTF8_Solid(font,text,clrFg);
	SDL_Rect rcDest = {x,y,0,0};
	SDL_SetAlpha(sText, SDL_SRCALPHA, alpha);
	SDL_BlitSurface(sText,NULL, screen,&rcDest);
	SDL_FreeSurface(sText);
	if (updateOn)
	  SDL_UpdateRect(screen,0,0,0,0);
	TTF_CloseFont(font);
	return 0;
}
//------------------------------------------------------------------
static void clearscreen () {
  SDL_FillRect(screen, NULL, SDL_MapRGBA(screen->format, red, green, blue, alpha));
  if (updateOn)
    SDL_UpdateRect (screen,0,0,0,0);
}
//------------------------------------------------------------
static int GrafikfensterAn (Uint16 breite, Uint16 hoehe) 
{   
    if ( fensteroffen == 1) return 1;  //Es kann nur ein Grafikfenster geben!

	//1. SDL muss initialisiert werden.
	if (SDL_Init(SDL_INIT_VIDEO | SDL_INIT_AUDIO) != 0) {
		//printf ("Kann SDL nicht initialisieren: %s\n", SDL_GetError ());
		return 1;
	}		
	//2. Bekanntmachung: Diese Funktion soll mit dem Programmende aufgerufen werden.
	// atexit (SDL_Quit);
	//3. Bildschirm: Hier kann man auch SDL_DOUBLEBUF sagen!
	screen = SDL_SetVideoMode (breite,hoehe, 32, SDL_DOUBLEBUF); //SDL_FULLSCREEN);
	if (screen == NULL) {
		//printf ("Bildschirm-Modus nicht setzbar: %s\n",SDL_GetError ());
		return 1;
	}
	SDL_WM_SetCaption( "LWB FU-Berlin: GO-Grafikfenster", 0);
	
	TTF_Init();
	
	red = 255;
	green = 255;
	blue = 255;
	alpha = 255; 
	clearscreen ();
	red   = 0;
	green = 0;
	blue  = 0;
	
	//Archiv-Surface erstellen
	archiv = SDL_ConvertSurface (screen, screen->format, SDL_HWSURFACE);
	if (archiv == NULL) {
		//printf ("Archiv-Surface konnte nicht erzeugt werden!\n");
		return 1;
	}
	
    // Open the audio device. The sound driver will try to give us
    // the requested format, but it might not succeed. The 'obtained'
    // structure will be filled in with the actual format data. 
    desired.freq = 44100;	// desired output sample rate 
    desired.format = AUDIO_S16;	// request signed 16-bit samples 
    desired.samples = 8192;	// this is more or less discretionary 
    desired.channels = 2;	// ask for stereo 
    desired.callback = AudioCallback;
    desired.userdata = NULL;	// we don't need this 
    if (SDL_OpenAudio(&desired, &obtained) < 0) {
    	//printf("Unable to open audio device: %s\n", SDL_GetError());
	    return 1;
    }
    // Initialisiere die Liste der möglichen Sounds (keiner aktiv zu Beginn) 
    int i;
    for (i = 0; i < MAX_PLAYING_SOUNDS; i++) {
	playing[i].active = 0;
    }

    // SDL's audio is initially paused. Start it. 
    SDL_PauseAudio(0);

	fensteroffen = 1;
	fensterzu = 0; 

	//Jetzt kommt die Event-Loop
	
	while (fensteroffen == 1 && SDL_WaitEvent(&event) != 0 ) {
		switch (event.type) {
			case SDL_KEYDOWN:
				if (tasteLesen)
				{
					gedrueckt = 1;                //Taste ist gerade heruntergedrückt.
					taste = event.key.keysym.sym;  //Das ist der Code der Taste auf der Tastatur.
					tiefe = event.key.keysym.mod;  //Gleichzeitig Steuerungstaste(n) gedrückt??
					//printf("%i,%i,%i\n",taste, gedrueckt, tiefe);
					tasteLesen = 0;
				}
				if (tastaturpuffer)
				{
					if (t_pufferende + 1 != t_pufferkopf)
					{
						t_puffer[t_pufferende] = ((Uint32) event.key.keysym.sym)*256*256 + (Uint32) 256*256*256*128 + ((Uint32) event.key.keysym.mod);
						t_pufferende++; //Umschlag auf 0 automatisch, da Uint8
					} 
				}
				break;
			case SDL_KEYUP:
				if (tasteLesen)
				{
					gedrueckt = 0; //Taste wurde gerade losgelassen.
					taste = event.key.keysym.sym;
					tiefe = event.key.keysym.mod;  //Gleichzeitig Steuerungstaste(n) gedrückt??
					//printf("%i,%i,%i\n",taste, gedrueckt, tiefe);
					tasteLesen = 0;
				}
				if (tastaturpuffer)
				{
					if (t_pufferende + 1 != t_pufferkopf)
					{
						t_puffer[t_pufferende] = ((Uint32) event.key.keysym.sym)*256*256 + ((Uint32) event.key.keysym.mod);
						t_pufferende++; //Umschlag auf 0 automatisch, da Uint8
					} 
				}
				break;
			case SDL_MOUSEMOTION:
				if (mausLesen)
				{   //BEi MOUSEMOTION GIBT ES NUR 3 MÖGLICHKEITEN FÜR EINE GEDRÜCKT-GEHALTENE TASTE: 1,2 oder 3
					// Dummerweise ist bei 3 der Tastenwert 4, daher Korrektur:
					mausTaste = (Uint8) event.button.button;
					if (mausTaste == 4)
						mausTaste--;
					mausX     = (Uint16) event.motion.x;
					mausY     = (Uint16) event.motion.y;
					mausLesen = 0;
				}
				if (mauspuffer)
				{   
					mausTaste = (Uint8) event.button.button;
					if (mausTaste == 4)
						mausTaste--;
					mausX     = (Uint16) event.motion.x;
					mausY     = (Uint16) event.motion.y;
					if (m_pufferende + 1 != m_pufferkopf)
					{
						m_puffer[m_pufferende] = ((Uint32) mausTaste)*256*256*256 + (((Uint32) mausX) <<12) + (Uint32) mausY;
						m_pufferende++;
					}
				}
				break;
			case SDL_MOUSEBUTTONDOWN:
				if (mausLesen)
				{
					mausTaste = (Uint8) event.button.button + 128; //+128: "pressed"
					mausX = (Uint16) event.motion.x;
					mausY = (Uint16) event.motion.y;
					mausLesen = 0;
				}
				if (mauspuffer)
				{
					mausTaste = (Uint8) event.button.button + 128; //+128: "pressed"
					mausX = (Uint16) event.motion.x;
					mausY = (Uint16) event.motion.y;
					if (m_pufferende + 1 != m_pufferkopf)
					{
						m_puffer[m_pufferende] = ((Uint32) mausTaste)*256*256*256 + (((Uint32) mausX) <<12) + (Uint32) mausY;
						m_pufferende++;
					}
				}
				break;
			case SDL_MOUSEBUTTONUP:
				if (mausLesen)
				{
					mausTaste = (Uint8) event.button.button + 64; //+64: "released"
					mausX = (Uint16) event.motion.x;
					mausY = (Uint16) event.motion.y;
					mausLesen = 0;
				}
				if (mauspuffer)
				{
					mausTaste = (Uint8) event.button.button + 64; //+64: "released"
					mausX = (Uint16) event.motion.x;
					mausY = (Uint16) event.motion.y;
					if (m_pufferende + 1 != m_pufferkopf)
					{
						m_puffer[m_pufferende] = ((Uint32) mausTaste)*256*256*256 + (((Uint32) mausX) <<12) + (Uint32) mausY;
						m_pufferende++;
					} 
				}
				break;
			case SDL_QUIT:
				//printf("Das Grafikfenster wurde geschlossen. Bye.\n");
				fensteroffen = 0;
				break;
		}
		
	}
	// Hier kommt man nur an, wenn das x angeklickt wurde.  
	// Die event-Loop wurde beendet, also soll das Fenster geschlossen werden!
	// Neue "Zeichenops" können nicht mehr starten, da fensteroffen==0 gilt.
	// Es muss aber HIER sichergestellt werden, dass keine "Zeichenoperation" mehr läuft,
	// damit es kein Segmentation-Fault gibt!
	
	// "OP_Lock ()"
	
	// Nun wird das SQL-Fenster komplett geschlossen:
	TTF_Quit ();
    // Pause and lock the sound system so we can safely delete our sound data. 
    SDL_PauseAudio(1);
    SDL_LockAudio();
    // Free our sounds before we exit, just to be safe.
    for (i=0; i < MAX_PLAYING_SOUNDS;i++) {
		if (playing[i].active ==1) {
			free(s[i].samples);
		}
	}
    // At this point the output is paused and we know for certain that the
    // callback is not active, so we can safely unlock the audio system. 
    SDL_UnlockAudio();
	SDL_CloseAudio();
	SDL_Quit ();
	// screen wird automatisch wieder freigegeben 
	SDL_FreeSurface(archiv); archiv = NULL;
	if (clipboard != NULL)
	{
		SDL_FreeSurface(clipboard); clipboard = NULL;
	}
	fensterzu = 1; 
	
	// "OP_Unlock ()"
	 
	return 0;
}
//------------------------------------------------------------------
static void Fenstertitel (const char *titel)
{
  SDL_WM_SetCaption(titel,0);
}
//------------------------------------------------------------------
static Uint8 FensterOffen ()
{
  return fensteroffen;
}
//------------------------------------------------------------------
static Uint8 FensterZu ()
{
  return fensterzu;
}
//-------------------------------------------------------------------
static void GrafikfensterAus ()
{
	SDL_Event user_event;
	user_event.type=SDL_QUIT;
	SDL_PushEvent(&user_event);
}
//-------------------------------------------------------------------
static void updateAus ()
{
	updateOn = 0;
}
//-------------------------------------------------------------------
static void updateAn ()
{
	updateOn = 1;
	SDL_Flip (screen);
}
//-------------------------------------------------------------------
static void zeichnePunkt (Sint16 x, Sint16 y)
{
  pixelRGBA (screen, x, y ,red, green,blue,alpha);
  if (updateOn) 
	SDL_UpdateRect (screen, x, y, 1, 1);
}
//--------------------------------------------------------------
static Uint32 gibPixel(Sint16 x, Sint16 y)
{
    int bpp = screen->format->BytesPerPixel;
    // Here p is the address to the pixel we want to retrieve 
    Uint8 *p = (Uint8 *)screen->pixels + y * screen->pitch + x * bpp;

    switch(bpp) {
    case 1:
        return *p;
        break;
    case 2:
        return *(Uint16 *)p;
        break;
    case 3:
        if(SDL_BYTEORDER == SDL_BIG_ENDIAN)
            return p[0] << 16 | p[1] << 8 | p[2];
        else
            return p[0] | p[1] << 8 | p[2] << 16;
        break;
    case 4:
        return *(Uint32 *)p;
        break;
    default:
        return 0;       // shouldn't happen, but avoids warnings 
    }
}
//--------------------------------------------------------------
static void zeichneKreis (Sint16 x, Sint16 y, Sint16 r, Uint8 full)
{ 
	if (full)
		filledCircleRGBA(screen,x,y,r,red,green,blue,alpha);
	else
		circleRGBA (screen, x,y,r,red, green, blue,alpha);
	if (updateOn)
		SDL_UpdateRect (screen,x-r,y-r,2*r+1,2*r+1);
}
//---------------------------------------------------------------
static void zeichneEllipse (Sint16 x, Sint16 y, Sint16 rx, Sint16 ry, Uint8 filled)
{	
	if (filled)
		filledEllipseRGBA (screen, x, y, rx, ry, red, green, blue, alpha);
	else
		ellipseRGBA (screen, x, y, rx, ry, red, green,blue, alpha);
	if (updateOn)
		SDL_UpdateRect (screen, x-rx, y-ry,2*rx+1,2*ry+1);
}
//---------------------------------------------------------------
static void stiftfarbe (Uint8 r, Uint8 g, Uint8 b)
{
    red = r;
    green = g;
    blue = b;
}
//---------------------------------------------------------------
static void transparenz (Uint8 t)
{
	alpha = t;
}
//---------------------------------------------------------------
static void zeichneStrecke (Sint16 x1, Sint16 y1, Sint16 x2, Sint16 y2)
{ 
	int upx,upy,breite,hoehe;
	lineRGBA (screen, x1,y1,x2,y2, red, green, blue, alpha);
	if (x1 <= x2)
	{
		upx    = x1;
		breite = x2 - x1 + 1;
	} 
	else
	{
		upx    = x2;
		breite = x1 - x2 + 1;
	}
	if (y1 <= y2)
	{
		upy   = y1;
		hoehe = y2 - y1 + 1;
	}
	else
	{
		upy   = y2;
		hoehe = y1 - y2 + 1;
	}
	if (updateOn)
		SDL_UpdateRect (screen,upx,upy,breite,hoehe);
}
//--------------------------------------------------------------
static void rechteck (Sint16 x1, Sint16 y1, Sint16 b, Sint16 h, Uint8 filled)
{	
	if (filled)
		boxRGBA (screen, x1, y1 ,x1+b-1, y1+h-1, red, green, blue, alpha);
	else
		rectangleRGBA (screen, x1, y1 , x1+b-1, y1+h-1, red, green,blue,alpha);
	if (updateOn)
		SDL_UpdateRect (screen, x1, y1, b, h);
}	
//--------------------------------------------------------------
static void kreissektor (Sint16 x, Sint16 y, Sint16 r, Sint16 w1, Sint16 w2, Uint8 filled)
{
	if (filled)
		filledPieRGBA (screen, x, y , r, w1, w2, red, green, blue, alpha);
	else
		pieRGBA (screen, x, y , r, w1, w2, red, green, blue, alpha);
	if (updateOn)
		SDL_UpdateRect (screen, x-r, y-r, 2*r+1, 2*r+1);
}
//---------------------------------------------------------------
Sint16 minimum (Sint16 x, Sint16 y, Sint16 z)
{
  if ((x <= y) && (x <=z))
    return x;
  else if ((y<=x) && (y<=z))
    return y;
  else
    return z;
}
//-------------------------------------------------------------
Sint16 maximum (Sint16 x, Sint16 y, Sint16 z)
{
  if ((x >= y) && (x >=z))
    return x;
  else if ((y>=x) && (y>=z))
    return y;
  else
    return z;
}
//---------------------------------------------------------------
static void dreieck (Sint16 x1, Sint16 y1, Sint16 x2, Sint16 y2, Sint16 x3, Sint16 y3, Uint8 filled)
{
	int upx,upy,breite,hoehe;

	upx = minimum (x1, x2, x3);
	upy = minimum (y1, y2, y3);
	breite = maximum (x1, x2, x3) - upx + 1;
	hoehe  = maximum (y1, y2, y3) - upy + 1;
	if (filled)
		filledTrigonRGBA(screen, x1,y1,x2,y2,x3,y3,red,green,blue,alpha);
	else
		trigonRGBA (screen, x1,y1,x2,y2,x3,y3,red,green,blue,alpha);
	if (updateOn)
		SDL_UpdateRect (screen, upx,upy,breite,hoehe);
}
//----------------------------------------------------------------
static void setcolorkey (Uint8 r, Uint8 g, Uint8 b, Uint8 key)
{
  ck_red = r;
  ck_green = g;
  ck_blue = b;
  colorkey = key;
}
//-----------------------------------------------------------------
static void ladeBild (Sint16 x, Sint16 y, char *cs)
{
	SDL_Surface *image;
	SDL_Rect src, dest;
	
	image = SDL_LoadBMP(cs);
	//printf ("Dateiname: %s\n",cs);
	if (image == NULL) {
		//printf("Bild konnte nicht geladen werden!\n");
		return;
	}
	src.x = 0;
	src.y = 0;
	src.w = image->w;
	src.h = image->h;
	
	dest.x = x;
	dest.y = y;
	dest.w = image->w;
	dest.h = image->h;
	
	if (colorkey)
		SDL_SetColorKey(image, SDL_SRCCOLORKEY, SDL_MapRGBA(screen->format, ck_red, ck_green, ck_blue, alpha));
	SDL_SetAlpha(image, SDL_SRCALPHA, alpha);
	SDL_BlitSurface(image, &src, screen, &dest);
	SDL_FreeSurface (image);  
	if (updateOn)
		SDL_UpdateRect(screen, x, y, dest.w, dest.h);
}
//---------------------------------------------------------
static void schreibe (Sint16 x, Sint16 y, char *cs)
{
	gfxPrimitivesSetFont(NULL, 0 ,0);
	stringRGBA (screen, x,y,cs,red, green, blue, alpha);
	if (updateOn)
		SDL_UpdateRect (screen,0,0,0,0);
}
//---------------------------------------------------------
static void ladeBildInsClipboard (char *cs)
{
	SDL_Surface *image;
	image = SDL_LoadBMP(cs);
	//printf ("Dateiname: %s\n",cs);
	if (image == NULL) {
		// printf("Bild konnte nicht geladen werden!\n");
		return;
	}
	SDL_FreeSurface (clipboard); //altes Clipboard freigeben
	clipboard = SDL_DisplayFormat (image);
	SDL_FreeSurface (image);
}
//---------------------------------------------------------
static void clipboardKopieren (Sint16 x, Sint16 y, Uint16 b, Uint16 h) 
{
	SDL_Surface *image;
	SDL_Rect src, dest;
	Uint32 rmask, gmask, bmask, amask;
	if (clipboard != NULL)
		SDL_FreeSurface (clipboard);
	#if SDL_BYTEORDER == SDL_BIG_ENDIAN
		rmask = 0xff000000;
		gmask = 0x00ff0000;
		bmask = 0x0000ff00;
		amask = 0x000000ff;
	#else
		rmask = 0x000000ff;
		gmask = 0x0000ff00;
		bmask = 0x00ff0000;
		amask = 0xff000000;
	#endif
	image = SDL_CreateRGBSurface(SDL_HWSURFACE | SDL_SRCALPHA, (int) b, (int) h, 32, rmask, gmask, bmask, amask);
	if (image == NULL) {
		// printf("Neues Clipboard konnte nicht erzeugt werden!\n");
		return;
	}
	src.x = x;
	src.y = y;
	src.w = b;
	src.h = h;
	dest.x = 0;
	dest.y = 0;
	dest.w = b;
	dest.h = h; 
	SDL_BlitSurface(screen, &src, image, &dest);
	SDL_UpdateRect (image, 0, 0, 0, 0);
	clipboard = SDL_DisplayFormat (image);
	SDL_FreeSurface (image);
}
//---------------------------------------------------------
static void clipboardEinfuegen (Sint16 x, Sint16 y)
{
	SDL_Rect src, dest;
	src.x = 0;
	src.y = 0;
	src.w = clipboard->w;
	src.h = clipboard->h;
	dest.x = x;
	dest.y = y;
	dest.w = clipboard->w;
	dest.h = clipboard->h;
	if (colorkey)
	{
		SDL_SetColorKey(clipboard, SDL_SRCCOLORKEY, SDL_MapRGBA(screen->format, ck_red, ck_green, ck_blue, alpha));
	} else {
		SDL_SetColorKey(clipboard, 0, SDL_MapRGBA(screen->format, ck_red, ck_green, ck_blue, alpha));
	}
	SDL_SetAlpha(clipboard, SDL_SRCALPHA, alpha); 
	SDL_BlitSurface(clipboard, &src, screen, &dest);
	if (updateOn)
		SDL_UpdateRect (screen, x, y, dest.w, dest.h);
}
//---------------------------------------------------------
static void archivieren ()
{
	SDL_Rect src, dest;
	src.x = 0;
	src.y = 0;
	src.w = screen->w;
	src.h = screen->h;
	dest = src;
	SDL_BlitSurface(screen, &src, archiv, &dest);
	SDL_UpdateRect(archiv, 0,0,0,0);
}
//----------------------------------------------------------
static void restaurieren (Sint16 x, Sint16 y, Uint16 b, Uint16 h)
{
	SDL_Rect src, dest;
	src.x = x;
	src.y = y;
	src.w = b;
	src.h = h;
	dest = src;
	SDL_BlitSurface(archiv, &src, screen, &dest);
	if (updateOn)
		SDL_UpdateRect (screen, x, y, b, h);
}
//---------------------------------------------------------------
static Uint32 tastaturLesen1 () 
{
	tasteLesen = 1;
	while (tasteLesen && fensteroffen) 
	{
		SDL_Delay (5);
	}
	return ((Uint32) taste)*256*256 + ((Uint32) gedrueckt)*256*256*256*128+ ((Uint32) tiefe);
}
//-------------------------------------------------------------
static void tastaturpufferAn () {
	t_pufferkopf = 0;
	t_pufferende = 0;
	tastaturpuffer = 1;
}
//-------------------------------------------------------------
static void tastaturpufferAus () {
	tastaturpuffer = 0;
}
//-------------------------------------------------------------
static Uint32 tastaturpufferLesen1 ()
{
	Uint32 erg;
	while (t_pufferende == t_pufferkopf && fensteroffen)
	{
		SDL_Delay (5);
	}
	erg = t_puffer[t_pufferkopf];
	t_pufferkopf++; //Überlauf von 255 auf 0 automatisch, da Uint8
	return erg;
}
//-------------------------------------------------------------
static Uint32 mausLesen1 ()
{
	mausLesen = 1;
	while (mausLesen && fensteroffen) 
	{
		SDL_Delay (5);
	}
	return ((Uint32) mausTaste)*256*256*256 + (((Uint32) mausX) << 12) + ((Uint32) mausY);
}
//--------------------------------------------------------------
static void mauspufferAn () {
	m_pufferkopf = 0;
	m_pufferende = 0;
	mauspuffer = 1;
}
//-------------------------------------------------------------
static void mauspufferAus () {
	mauspuffer = 0;
}
//-------------------------------------------------------------
static Uint32 mauspufferLesen1 ()
{
	Uint32 erg;
	while (m_pufferende == m_pufferkopf && fensteroffen)
	{
		SDL_Delay (5);
	}
	erg = m_puffer[m_pufferkopf];
	m_pufferkopf++; //Überlauf von 255 auf 0 automatisch, da Uint8
	return erg;
}
//-------------------------------------------------------------
*/
import "C"

import ( "time" ; "unsafe" )

var grafikschloss   = make (chan int,1)
var tastaturschloss = make (chan int,1)
var fensterschloss  = make (chan int,1)
var mausschloss     = make (chan int,1)
var fensterbreite,fensterhoehe uint16

// Es gibt 4 Tastenbelegungen: Standard, SHIFT, ALT GR, ALT GR mit SHIFT.

var z1 [4]string = [4]string{ ",-.", ";_:",  "·–…", "×—÷"}
var z2 [4]string = [4]string{"0123456789", "=!\"§$%&/()", "}¹²³¼½¬{[]", "°¡⅛£¤⅜⅝⅞™±"}
var z3 [4]string = [4]string{"abcdefghijklmnopqrstuvwxyz", "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "æ“¢ð€đŋħ→̣ĸłµ”øþ@¶ſŧ↓„ł«»←",  "Æ‘©Ð€ªŊĦı˙&Łº’ØÞΩ®ẞŦ↑‚Ł‹›¥"}
var z4 [4]string = [4]string{",/*-+",",/*-+",",/*-+",",/*-+"} //NUM-Block
var z5 [4]string = [4]string{"0123456789","0123456789","0123456789","0123456789"} //NUM-BLOCK
var taste_belegung [4][320]rune //vier Belegungen pro Taste
//-----------------------------------------------------------------------------

func lock () { grafikschloss <- 1 }
func unlock () { <- grafikschloss }
func t_lock () { tastaturschloss <- 1 }
func t_unlock () { <- tastaturschloss }
func m_lock () { mausschloss <- 1}
func m_unlock () { <- mausschloss }

// Vor.: Das Grafikfenster ist nicht offen. Es gilt: breite <=1920, hoehe <=1200.
//
// Eff.: Ein Fenster mit einer 'Zeichenfläche' von breite x hoehe Pixeln wurde geöffnet
// Die Zeichenfarbe ist Schwarz, ohne Transparenz. Der Ursprung (0,0) ist oben links
// im Fenster. Die x-Koordinate wächst horizontal nach rechts, die y-Koordinate vertikal nach unten.
func Fenster (breite, hoehe uint16) {
	lock ()
	if fensterZu () {
		if breite > 1920 {breite = 1920}
		if hoehe > 1200 {hoehe = 1200}
		fensterhoehe = hoehe
		fensterbreite = breite
		go C.GrafikfensterAn (C.Uint16(breite), C.Uint16(hoehe))
		for !FensterOffen () {
			time.Sleep (100 * 1000 * 1000) //Unter Windows notwendig!!
		}
	}
	unlock ()
}

// Vor.: -
//
// Erg.:  True ist geliefert, gdw. das Grafikfenster offen ist.
func FensterOffen () bool {
	return uint8(C.FensterOffen ()) == 1
}

func fensterZu () bool { //interne Hilfsfunktion
	return uint8(C.FensterZu ()) == 1
}

// Vor.: Das Grafikfenster ist offen.
// 
// Eff.: Das Grafikfenster ist geschlossen.
func FensterAus () {
	lock ()
	if FensterOffen () {
		C.GrafikfensterAus ()
		for !fensterZu () {
			time.Sleep (100 * 1000 * 1000) 
		}
	}
	unlock ()
}

// Vor.: Das Grafikfenster ist offen.
//
// Erg.: Die Anzahl der Grafikfensterzeilen (Pixelzeilen) des gfx-Fensters
// ist geliefert.
func Grafikzeilen () uint16 {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	return fensterhoehe
}

// Vor.: Das Grafikfenster ist offen.
//
// Erg.: Die Anzahl der Grafikfensterspalten (Pixelspalten) des gfx-Fensters
// ist geliefert.
func Grafikspalten () uint16 {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	return fensterbreite
}

// Vor.: Das Grafikfenster ist offen.
//
// Eff.: Das gfx-Fenster hat sichtbar den neuen Fenstertitel s.
// In der Regel verwendet man hier den Programmnamen.
func Fenstertitel (s string) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	cs:=C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.Fenstertitel (cs)
	unlock ()
}

// Vor.: Das Grafikfenster ist offen.
//
// Eff.: Alle Pixel des Grafikfenster haben nun die aktuelle Stiftfarbe,
// d.h., der Inhalt des Fensters ist gelöscht.
func Cls () {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.clearscreen ()
	unlock ()
}

// Vor.: Das Grafikfenster ist offen.
//
// Eff.: Die Zeichenfarbe ist gemäß dem RGB-Farbmodell neu gesetzt.
// Beispiel: Stiftfarbe (0xFF, 0, 0) ist Rot. Die Transparenz der Stiftfarbe
// kann mit der Funktion Transparenz eingestellt werden.
func Stiftfarbe (r,g,b uint8) {
	lock ()
	C.stiftfarbe (C.Uint8 (r), C.Uint8 (g), C.Uint8 (b))
	unlock ()
}

// Vor.: Das Grafikfenster ist offen.
//
// Eff.: Die Transparenz der Stiftfarbe bzw. die von "Zeichenoperationen" ist neu gesetzt.
// 0 bedeutet keine Transparenz (Standard), 255 komplett durchsichtig. Wenn also etwas
// nach dem Aufruf gezeichnet wird, so scheint vorher Gezeichnetes ggf. durch. 
func Transparenz (t uint8) {
	lock ()
	C.transparenz (C.Uint8(255-t))
	unlock ()
}

// Vor.: Das Grafikfenster ist offen.
//
// Eff.: An der Position (x,y) ist  ein Punkt in der aktuellen Stiftfarbe gesetzt.
func Punkt (x,y uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.zeichnePunkt (C.Sint16(x), C.Sint16(y))
	unlock ()
}

// Vor.: Das Grafikfenster ist offen. 
//
// Erg.: Der Rot-, Grün- und Blauanteil des Punktes mit den Koordinaten
// (x,y) im Grafikfenster ist geliefert.
func GibPunktfarbe (x,y uint16) (r,g,b uint8) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	pixel:=uint32 (C.gibPixel(C.Sint16(x),C.Sint16(y)))
	r = uint8(pixel >> 16)
	g = uint8(pixel >> 8)
	b = uint8(pixel)
	unlock ()
	return
}

// Vor.: Das Grafikfenster ist offen.
//
// Eff.: Von der Position (x1,y1) bis (x2,y2) ist eine Strecke mit der 
// Strichbreite 1 Pixel in der aktuellen Stiftfarbe gezeichnet.
func Linie (x1,y1,x2,y2 uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.zeichneStrecke (C.Sint16(x1),C.Sint16(y1),C.Sint16(x2),C.Sint16(y2))
	unlock ()
}

// Vor.: Das Grafikfenster ist offen.
//
// Eff.: Um den Mittelpunkt M (x,y) ist ein Kreis mit dem Radius r mit der 
// Strichbreite 1 Pixel in der aktuellen Stiftfarbe gezeichnet.
func Kreis (x,y,r uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.zeichneKreis(C.Sint16(x),C.Sint16(y),C.Sint16(r),0)
	unlock ()
}

// Vor.: Das Grafikfenster ist offen.
//
// Eff.: Um den Mittelpunkt M (x,y) ist ein ausgefüllter Kreis mit dem 
// Radius r in der aktuellen Stiftfarbe gezeichnet.
func Vollkreis (x,y,r uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.zeichneKreis(C.Sint16(x),C.Sint16(y),C.Sint16(r),1)
	unlock ()
}

// Vor.: Das Grafikfenster ist offen.
//
// Eff.: Um den Mittelpunkt M (x,y) ist mit der horizontalen Halbachse rx 
// und der vertikalen Halbachse ry mit der Strichbreite 1 Pixel in
// der aktuellen Stiftfarbe eine Ellipse gezeichnet.
func Ellipse (x,y,rx,ry uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.zeichneEllipse(C.Sint16(x),C.Sint16(y),C.Sint16(rx),C.Sint16(ry),0)
	unlock ()
}

// Vor.: Das Grafikfenster ist offen.
//
// Eff.: Um den Mittelpunkt M (x,y) ist mit der horizontalen Halbachse rx
// und der vertikalen Halbachse ry in der aktuellen Stiftfarbe eine
// ausgefüllte Ellipse gezeichnet.
func Vollellipse (x,y,rx,ry uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.zeichneEllipse(C.Sint16(x),C.Sint16(y),C.Sint16(rx),C.Sint16(ry),1)
	unlock ()
}

// Vor.: Das Grafikfenster ist offen.
//
// Eff.: Um den Mittelpunkt M (x,y) ist mit dem Radius r in der aktuellen
// Stiftfarbe ein Kreisektor(Tortenstück:-)) gezeichnet. w1 ist 
// dabei der Startwinkel in Grad, w2 der Endwinkel in Grad. Ein 
// Winkelmaß von 0 Grad bedeutet in Richtung Osten geht es los, dann
// entgegengesetzt zum Uhrzeigersinn.
func Kreissektor (x,y,r,w1,w2 uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.kreissektor (C.Sint16(x),C.Sint16(y),C.Sint16(r),360-C.Sint16(w2),360-C.Sint16(w1),0)
	unlock ()
}

// Vor.: Das Grafikfenster ist offen.
//
// Eff.: Um den Mittelpunkt M (x,y) ist mit dem Radius r  in der aktuellen
// Stiftfarbe ein gefüllter Kreissegment gezeichnet. w1 ist dabei 
// der Startwinkel in Grad, w2 der Endwinkel in Grad. Ein Winkelmaß
// von 0 Grad bedeutet in Richtung Osten geht es los, dann entgegen- 
// gesetzt zum Uhrzeigersinn.
func Vollkreissektor (x,y,r,w1,w2 uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.kreissektor (C.Sint16(x),C.Sint16(y),C.Sint16(r),360-C.Sint16(w2),360-C.Sint16(w1),1)
	unlock ()
}

// Vor.: Das Grafikfenster ist offen.
//
// Eff.: In der aktuellen Stiftfarbe ist ein Rechteck gezeichnet. Die 
// Position (x1,y1) gibt die linke obere Ecke des Rechtecks an, b 
// die Breite in x-Richtung, h die Höhe in y-Richtung. Die Seiten
// des Rechtecks verlaufen parallel zu den Achsen.
func Rechteck (x1,y1,b,h uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.rechteck (C.Sint16(x1),C.Sint16(y1),C.Sint16(b),C.Sint16(h),0)
	unlock ()
}

// Vor.: Das Grafikfenster ist offen.
//
// Eff.: In der aktuellen Stiftfarbe ist ein gefülltes Rechteck gezeichnet.
// Die Position (x1,y1) gibt die linke obere Ecke des Rechtecks an,
// b die Breite in x-Richtung, h die Höhe in y-Richtung. Die Seiten
// des Rechtecks verlaufen parallel zu den Achsen.
func Vollrechteck (x1,y1,b,h uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.rechteck (C.Sint16(x1),C.Sint16(y1),C.Sint16(b),C.Sint16(h),1)
	unlock ()
}

// Vor.: Das Grafikfenster ist offen.
//
// Eff.: In der aktuellen Stiftfarbe ist ein Dreieck mit den Eckpunkt-
// koordinaten (x1,y1), (x2,y2) und (x3,y3) gezeichnet.
func Dreieck (x1,y1,x2,y2,x3,y3 uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.dreieck (C.Sint16(x1),C.Sint16(y1),C.Sint16(x2),C.Sint16(y2),C.Sint16(x3),C.Sint16(y3),0)
	unlock ()
}

// Vor.: Das Grafikfenster ist offen.
//
// Eff.: In der aktuellen Stiftfarbe ein gefülltes Dreieck mit den
// Eckpunktkoordinaten (x1,y1), (x2,y2) und (x3,y3) gezeichnet.
func Volldreieck (x1,y1,x2,y2,x3,y3 uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.dreieck (C.Sint16(x1),C.Sint16(y1),C.Sint16(x2),C.Sint16(y2),C.Sint16(x3),C.Sint16(y3),1)
	unlock ()
}

// Vor.: Das Grafikfenster ist offen. s beinhaltet maximal 255 Bytes und
// ist ein ASCII-Code-String.
//
// Eff.: In der aktuellen Stiftfarbe ist der Text s hingeschrieben ohne 
// den umgebenden Hintergrund zu verändern. Die Position (x,y) ist
// die linke obere Ecke des Bereichs des ersten Buchstaben von S. 
func Schreibe (x,y uint16, s string) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	cs:= C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.schreibe (C.Sint16(x), C.Sint16 (y), cs)
	unlock ()
}
	
// Vor.: s gibt die ttf-Datei des Fonts mit vollständigem Pfad an.
// groesse gibt die gewünschte Punkthöhe der Buchstaben an.
//
// Eff.: Wenn es die ttf-Datei gibt, so ist der angegebene Font nun der
// aktuelle Font, der bei Aufruf von SchreibeFont () verwendet wird.
//
// Erg.: -true- ist geliefert, gdw. die ttf-Datei an der Stelle lag und 
// der Font als aktueller Font gesetzt werden konnte.
func SetzeFont (s string, groesse int) (erg bool) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	cs:=C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	if int(C.setFont(cs,C.int(groesse)))==0 {
		erg = true 
	} else {
		erg = false
	}
	unlock()
	return
}

// Vor.: keine
//
// Erg.: Der mit SetzeFont () hinterlegte Pfad inklusive Dateiname
// des aktuell gewünschten Fonts ist geliefert.
func GibFont () (erg string) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	cs:=C.getFont ()
	erg = C.GoString(cs)
	// C.free(unsafe.Pointer(cs)) knallt mit 'invalid pointer' -- Warum?? Unklar! SC
	unlock()
	return
}

// Vor.: Das Grafikfenster ist offen. s beinhaltet maximal 255 Bytes.
//
// Eff.: In der aktuellen Stiftfarbe ist der Text s mit dem zuletzt mit
// SetzeFont() gesetzten Font hingeschrieben ohne 
// den Hintergrund zu verändern. Die Position (x,y) ist die linke
// obere Ecke des Bereichs des ersten Buchstaben von S. 
func SchreibeFont (x,y uint16, s string) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	if len (s) > 0 { // Bei einem Leerstring knallt 'C.write' unter Linux mit Segmentation Fault! Erst seit 2020! Unklar! SC
		cs:=C.CString(s)
		if int(C.write (C.Sint16(x),C.Sint16(y),cs)) == 1 {
			println ("FEHLER: Kein aktueller Font: ", C.GoString(C.getFont()))
		}
		C.free(unsafe.Pointer(cs))
	}
	unlock()
}

// Vor.: Das Grafikfenster ist offen. s beinhaltet maximal 255 Bytes und
// stellt den Dateinamen eines Bildes im bmp-Format dar.
//
// Eff.: Ab der Position (x,y) ist das angegebene rechteckige Bild gemäß
// der aktuell eingestellten Transparenz eingefügt. Die Position ist 
// die linke obere Ecke des Bildes. Die Bildkanten verlaufen parallel 
// zu den Achsen.
func LadeBild (x,y uint16, s string) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	cs:= C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.setcolorkey (C.Uint8(0),C.Uint8(0),C.Uint8(0),C.Uint8(0))
	C.ladeBild (C.Sint16(x),C.Sint16(y),cs)
	unlock ()
}

// Vor.: Das Grafikfenster ist offen. s beinhaltet maximal 255 Bytes und
// stellt den Dateinamen eines Bildes im bmp-Format dar.
// r,g und b geben eine Pixelfarbe an (ColorKey).
//
// Eff.: Ab der Position (x,y) ist das angegebene rechteckige Bild gemäß
// der eingestellten Transparenz eingefügt. Die Position ist die linke obere
// Ecke des Bildes. Die Bildkanten verlaufen parallel zu den Achsen. Alle
// Pixel des Bildes mit den Farbwerten r,g und b werden jedoch vollkommen
// transparent dargestellt! Ursprüngliche Pixel im Grafikfenster werden
// hier nicht überzeichnet!
func LadeBildMitColorKey (x,y uint16, s string,r,g,b uint8) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	cs:= C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.setcolorkey (C.Uint8(r),C.Uint8(g),C.Uint8(b),C.Uint8(1))
	C.ladeBild (C.Sint16(x),C.Sint16(y),cs)
	unlock ()
}

// Vor.: Das Grafikfenster ist offen. s beinhaltet maximal 255 Bytes und
// stellt den Dateinamen eines Bildes im bmp-Format dar.
//
// Eff.: Das angegebene Bild ist in einen Zwischenspeicher (das Clipboard)
// geladen. Vorher im Clipboard enthaltene Daten wurden damit überschrieben.	
func LadeBildInsClipboard (s string) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	cs:= C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.ladeBildInsClipboard (cs)
	unlock ()
}

// Vor.: Das Grafikfenster ist offen.
//
// Eff.: Der gesamte Inhalt des Fensters ist in einen (versteckten)
// Zwischenspeicher kopiert. Daten, die vorher in diesem Zwischen- 
// speicher waren, wurden überschrieben. 
func Archivieren () {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.archivieren ()
	unlock ()
}

// Vor.: Das Grafikfenster ist offen. Archivieren wurde vorher mindestens
// einmal aufgerufen und seit dem das Fenster nicht geschlossen.
//
// Eff.: Der angegebene rechteckige Bereich des versteckten Zwischenspeichers
// (s. Archivieren) ist an seine ursprüngliche Stelle ins Grafikfenster
// zurückkopiert. Die gesetzte Transparenz hat keinen Einfluss auf die Funktion.
func Restaurieren (x1,y1,b,h uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.restaurieren (C.Sint16(x1),C.Sint16(y1),C.Uint16(b),C.Uint16(h))
	unlock ()
}

// Vor.: Das Grafikfenster ist offen.
//
// Eff.: Der angegebene rechteckige Grafikfensterbereich ist in einem 
// Zwischenspeicher (das Clipboard) kopiert. Daten, die vorher in
// diesem Zwischenspeicher waren, wurden überschrieben.
func Clipboard_kopieren (x,y,b,h uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.clipboardKopieren (C.Sint16(x), C.Sint16(y), C.Uint16(b), C.Uint16 (h))
	unlock ()
}

// Vor.: Das Grafikfenster ist offen, Clipboard_kopieren wurde vorher
// mindestens einmal aufgerufen und seitdem wurde das Fenster nicht geschlossen.
//
// Eff.: Der Inhalt des Zwischenspeichers (Clipboard) ist an die angege-
// bene Position (x,y) ins Grafikfenster kopiert. Dort vorher 
// vorhandene Daten wurden überschrieben, wobei die gesetzte Transparenz
// entsprechenden Einfluss hatte.
func Clipboard_einfuegen(x,y uint16) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.setcolorkey (C.Uint8(0),C.Uint8(0),C.Uint8(0),C.Uint8(0))
	C.clipboardEinfuegen(C.Sint16(x), C.Sint16(y))
	unlock ()
}

// Vor.: Das Grafikfenster ist offen, Clipboard_kopieren wurde vorher
// mindestens einmal aufgerufen und seitdem wurde das Fenster 
// nicht geschlossen. r,g und b geben eine Pixelfarbe an.
//
// Eff.: Der Inhalt des Zwischenspeichers (Clipboard) ist an die angege-
// bene Position (x,y) ins Grafikfenster unter Beachtung der gesetzten
// Transparenz kopiert. Alle Pixel des Clipboards mit dem durch r,g und b
// festgelegten Farbwertes sind jedoch vollkommen transparent und ändern
// so das ursprüngliche Pixel im Grafikfenster an dieser Stelle nicht. 
func Clipboard_einfuegenMitColorKey (x,y uint16, r,g,b uint8) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.setcolorkey (C.Uint8(r),C.Uint8(g),C.Uint8(b),C.Uint8(1))
	C.clipboardEinfuegen(C.Sint16(x), C.Sint16(y))
	unlock ()
}

// Vor.: Das Grafikfenster ist offen. Sperren wurde noch nicht aufgerufen bzw.
// der Aufruf wurde mit einem Aufruf von Entsperren 'neutralisiert'. 
//
// Eff.: Das Grafikfenster ist nun nur noch vom aufrufenden Prozess 'beschreibbar',
// wenn alle anderen Prozesse vor einem Schreibzugriff auf das Grafikfenster
// ebenfalls Sperren aufrufen. Gegebenenfalls war der aufrufende Prozess
// solange blockiert, bis er den Zugriff erhielt. Andere Prozesse, 
// die nun Sperren ausführen, sind blockiert.
func Sperren () {
	fensterschloss <- 1
}

// Vor.: Das Grafikfenster ist offen. Sperren wurde aufgerufen und seit 
// dem das Grafikfenster nicht geschlossen.
//
// Eff.: Das Grafikfenster ist für andere Prozesse wieder zum 'Beschreiben'
// freigegeben.
func Entsperren () {
	<- fensterschloss
}

// Vor.: Das Grafikfenster ist offen.
//
// Eff.: Die abgesetzten gfx-Anweisungen werden nicht sofort im Fenster,
// sondern lediglich im 'Double-Buffer-Bereich' verdeckt durchgeführt.
func UpdateAus () {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.updateAus ()
	unlock ()
}

// Vor.: Das Grafikfenster ist offen und wurde nach einem 'UpdateAus()'
// nicht geschlossen.
//
// Eff.: Alle nach 'UpdateAus ()' durchgeführten Änderungen durch abgesetzte 
// Grafikbefehle sind nun sichtbar geworden. Folgende gfx-Anweisungen
// werden wieder direkt umgesetzt.
func UpdateAn () {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	C.updateAn ()
	unlock ()
}

// Vor.: Das Grafikfenster ist offen.
//  
// Erg.: Der aufrufende Prozess war solange blockiert, bis eine Taste 
// auf der Tastatur gedrückt oder losgelassen wurde. Geliefert
// ist mit 'taste' die Tastennummer. 'gedrückt' ist 1 (0),falls die
// Taste gedrückt (losgelassen) wurde. 'tiefe' liefert die Kombination
// der gedrückten Steuerungstasten.
func TastaturLesen1 () (taste uint16, gedrueckt uint8, tiefe uint16) {
	var tastenwert uint32
	t_lock ()
	tastenwert = uint32(C.tastaturLesen1 ())
	t_unlock ()
	tiefe = uint16 (tastenwert % 65536)
	tastenwert = tastenwert >> 16
	gedrueckt = uint8(tastenwert >> 15)
	taste = uint16(tastenwert % 32768) //oberstes Bit rausschieben
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	return
}

// Vor.: Die Tastaturbelegung ist deutsch.
//
// Erg.: Wenn -tiefe- nur SHIFT oder STANDARD (also kein SHIFT) in Kombination
// mit NUMLOCK und/oder ALT GR ist und eine Tastaturzeichen-Taste
// mit -taste- übergeben wurde (also keine Steuertastenkombination),
// so ist das entsprechende Tastaturzeichen als Rune geliefert.
// Andernfalls ist rune(0) geliefert.
// -tiefe- und -taste- erhält man i.d.R. durch Tastaturlesen1().
func Tastaturzeichen (taste, tiefe uint16) rune {
	switch tiefe {
		case 0,4096, 8192+1, 8192+2, 8192+3,4096+8192+1,4096+8192+2,4096+8192+3:	// kein SHIFT, kein ALT GR, NUMLOCK an oder aus, CAPSLOCK an mit SHIFT
			return taste_belegung[0][taste]
		case 1,2,3,4096+1, 4096+2, 4096+3, 8192, 4096+8192:  // SHIFT, kein ALT GR, NUMLOCK an oder aus, CAPSLOCK an ohne SHIFT
			return taste_belegung[1][taste]
		case 16384, 16384 + 4096, 16384+8192+1,16384+8192+2,16384+8192+3,16384+8192+4096+1,16384+8192+4096+2,16384+8192+4096+3:  // kein SHIFT, ALT GR, NUMLOCK an oder aus, CAPSLOCK an mit SHIFT
			return taste_belegung[2][taste]
		case 16384+1, 16384+2, 16384+3, 16384+4096+1, 16384+4096+2, 16384+4096+3, 16384+8192, 16384+8192+4096: // ALT GR und SHIFT, NUMLOCK an oder aus, CAPSLOCK an ohne SHIFT
			return taste_belegung[3][taste]
		default:
		return 0
	}
}

// Vor.: Das Grafikfenster ist offen.
//  
// Eff.: Ab jetzt werden bis zu 255 Tastaturereignisse in einem 
// versteckten Tastaturpuffer zwischengespeichert. Darüber hinaus-
// gehende eingehende Tastaturevents gehen verloren.
func TastaturpufferAn () {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	t_lock ()
	C.tastaturpufferAn ()
	t_unlock ()
}

// Vor.: Das Grafikfenster ist offen.
//   
// Eff.: Der Tastaturpuffer ist aus. Enthaltene Events sind verloren.
func TastaturpufferAus (){
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	t_lock ()
	C.tastaturpufferAus ()
	t_unlock ()
}

// Vor.: Das Grafikfenster ist offen.
//  
// Erg.: Das vorderste Element (gespeicherte Event) des Tastaturpuffers 
// ist ausgelesen, aus dem Puffer entfernt  und zurückgegeben: Geliefert
// ist mit 'taste' die Tastennummer. 'gedrückt' ist 1 (0),falls die
// Taste gedrückt (losgelassen) wurde. 'tiefe' liefert die Kombination
// der gedrückten Steuerungstasten.
// War der Puffer leer, so war der aufrufende Prozess solange 
// blockiert, bis etwas gelesen werden konnte.
func TastaturpufferLesen1 () (taste uint16, gedrueckt uint8, tiefe uint16) {
	var tastenwert uint32
	t_lock ()
	tastenwert = uint32(C.tastaturpufferLesen1 ())
	t_unlock ()
	tiefe = uint16 (tastenwert % 65536)
	tastenwert = tastenwert >> 16
	gedrueckt = uint8(tastenwert >> 15)
	taste = uint16(tastenwert % 32768)
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	return
}

// Vor.: Das Grafikfenster ist offen.
//  
// Erg.: Der aufrufende Prozess war solange blockiert, bis Daten von der
// Maus gelesen werden konnten. Mit 'taste' erhält man die Nummer 
// der betreffenden Maustaste. Mit 'status' (1/0/-1), ob sie gedrückt
// bzw. unverändert ist oder losgelassen wurde. 'mausX' und 'mausY' 
// sind die Koordinaten der Mauszeigerspitze.
func MausLesen1 () (taste uint8,status int8, mausX, mausY uint16) {
	var tastenwert uint32
	m_lock ()
	tastenwert = uint32(C.mausLesen1 ())
	m_unlock ()
	taste = uint8 (tastenwert >> 24)
	if taste < 64 {
		status=0   //Zustand wird gehalten
	} else if taste > 128 {
		status = 1 //gerade gedrückt
		taste = taste - 128
	} else  {//zwischen 64 und 128
		status = -1 //gerade losgelassen
		taste = taste-64
	} 
	mausY = uint16 (tastenwert % 4096)
	tastenwert = tastenwert >> 12
	mausX = uint16 (tastenwert % 4096)
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	return
}

// Vor.: Das Grafikfenster ist offen.
//  
// Eff.: Ab jetzt werden bis zu 255 Mausereignisse (Events) zwischen-
// gespeichert.Darüber hinaus eingehende Maus-Events gehen verloren.
func MauspufferAn () {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	m_lock()
	C.mauspufferAn ()
	m_unlock ()
}

// Vor.: Das Grafikfenster ist offen. 
// 
// Eff.: Der Mauspuffer ist deaktiviert. Enthaltene Ereignisse sind 
// verloren.
func MauspufferAus (){
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	m_lock ()
	C.mauspufferAus ()
	m_unlock ()
}

// Vor.: Das Grafikfenster ist offen.
//
// Erg.: Das vorderste Mausereignis ist aus dem Puffer gelesen, dort 
// entfernt und zurückgegeben: Mit 'taste' erhält man die Nummer der betreffenden
// Maustaste. Mit 'status' (1/0/-1), ob sie gedrückt bzw. unverändert ist oder
// losgelassen wurde. 'mausX' und 'mausY' sind die Koordinaten der Mauszeigerspitze.
// War der Puffer leer, so war der aufrufende Prozess solange 
// blockiert, bis er etwas lesen konnte.
func MauspufferLesen1 () (taste uint8,status int8, mausX, mausY uint16) {
	var tastenwert uint32
	m_lock ()
	tastenwert = uint32(C.mauspufferLesen1 ())
	m_unlock ()
	taste = uint8 (tastenwert >> 24)
	if taste < 64 {
		status=0   //Zustand wird gehalten
	} else if taste > 128 {
		status = 1 //gerade gedrückt
		taste = taste - 128
	} else  {//zwischen 64 und 128
		status = -1 //gerade losgelassen
		taste = taste-64
	} 
	mausY = uint16 (tastenwert % 4096)
	tastenwert = tastenwert >> 12
	mausX = uint16 (tastenwert % 4096)
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	return
} 

// Vor.: Das Grafikfenster ist offen. s ist der Dateiname der wav-Datei
// inklusive Pfad. Zum Zeipunkt des Aufrufs werden gerade höchstens 9 .wav-Dateien abgespielt.
//
// Eff.: Die angegebene wav-Datei wird ab jetzt auch abgespielt.
// Das Programm läuft ohne Verzögerung weiter.
func SpieleSound(s string) {
	if !FensterOffen() { panic ("Das gfx-Fenster ist nicht offen!") }
	lock ()
	cs:=C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	erg:= int(C.LoadAndPlaySound(cs))
	if erg == 1 {
		println("Es werden schon die max. Anzahl an Sounds abgespielt!")
	}
	if erg == 2 {
		println ("Konnte Sounddatei nicht laden! --> ", s)
	}
	unlock()
}

// INTERN
// Vor.: data stellt die Bytefolge einer WAVE-Datei dar.
//       wartezeit ist die Abwartezeit nach dem Anspielen der WAV-Datei in ms.
// Eff.: Die 'WAV-Datei' wird bzw. ist gerade abgespielt. Der Programmablauf
//       ist dafür um wartezeit ms verzögert worden.

func spieleRAMWAV (data []byte,wartezeit uint32) {
	lock ()
	erg:= int (C.LadeUndSpieleNote(unsafe.Pointer(&data[0]),C.int(len(data))))
	if erg == 1 {
		println("Es werden schon die max. Anzahl an Sounds abgespielt!")
	}
	if erg == 2 {
		println ("Die Daten entsprechen keiner WAV-Datei! Daten nicht geladen!")
	}
	unlock()
	time.Sleep (time.Duration(int64(wartezeit)*1e6))
}

func init ()  {
	// Es folgt die Initalisierung der Tastaturbelegung auf Deutsch.
	// Das wird für die Funktion 'Tastaturzeichen(taste, tiefe) rune' benötigt.
	for i:=0; i < 4; i++ {
		index:= 0; for _,e:= range z1[i] {taste_belegung[i][index+44] = e; index++}
		index = 0; for _,e:= range z2[i] {taste_belegung[i][index+48] = e; index++}
		index = 0; for _,e:= range z5[i] {taste_belegung[i][index+256] = e; index++} //Num-Block
		index = 0; for _,e:= range z4[i] {taste_belegung[i][index+266] = e; index++} //Num-Block
		index = 0; for _,e:= range z3[i] {taste_belegung[i][index+97] = e; index++}
	}
	// kein SHIFT, kein ALT GR
	taste_belegung[0][43]='+'  ; taste_belegung[0][35] ='#' ; taste_belegung[0][252]='ü'
	taste_belegung[0][246]='ö' ; taste_belegung[0][228]='ä' ; taste_belegung[0][223]='ß'
	taste_belegung[0][180]='´' ; taste_belegung[0][94] ='^' ; taste_belegung[0][60] ='<'
	taste_belegung[0][32]=' '
	// SHIFT, kein ALT GR
	taste_belegung[1][43] ='*' ; taste_belegung[1][35] ='\'' ; taste_belegung[1][252]='Ü'
	taste_belegung[1][246]='Ö' ; taste_belegung[1][228]='Ä'  ; taste_belegung[1][223]='?'
	taste_belegung[1][180]='`' ; taste_belegung[1][94] ='°'  ; taste_belegung[1][60]='>'
	taste_belegung[1][32]=' '
	// kein SHIFT, ALT GR
	taste_belegung[2][43] ='~' ; taste_belegung[2][35] ='`' ; taste_belegung[2][252]='¨'
	taste_belegung[2][246]='˝' ; taste_belegung[2][228]='^' ; taste_belegung[2][223]='\\'
	taste_belegung[2][180]='¸' ; taste_belegung[2][94] ='¬' ; taste_belegung[2][60] ='|'
	taste_belegung[2][32] =' '
	// SHIFT, ALT GR
	taste_belegung[3][43] ='¯' ; taste_belegung[3][35] ='`' ; taste_belegung[3][252]='¨'
	taste_belegung[3][246]='˝' ; taste_belegung[3][228]='^' ; taste_belegung[3][223]='¿'
	taste_belegung[3][180]='¸' ; taste_belegung[3][94] ='¬' ; taste_belegung[3][60] ='¦'
	taste_belegung[3][32] =' '
}
//------------------------ENDE----------------------------------
