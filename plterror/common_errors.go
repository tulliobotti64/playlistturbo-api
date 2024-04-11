package plterror

import "net/http"

var ErrBadSyntax = &PLTError{message: "ERR_BAD_SYNTAX", status: http.StatusBadRequest}
var ErrServerError = &PLTError{message: "ERR_INTERNAL_SERVER_ERROR", status: http.StatusInternalServerError}
var ErrNoDataAdd = &PLTError{message: "ERR_NO_DATA_ADDED", status: http.StatusInternalServerError}
var InvalidExtension = &PLTError{message: "ERR_INVALID_EXTENSION", status: http.StatusBadRequest}
var InvalidSongPath = &PLTError{message: "ERR_INVALID_SONG_PATH", status: http.StatusBadRequest}
var InvalidGenre = &PLTError{message: "ERR_INVALID_GENRE", status: http.StatusBadRequest}
var Tabelavazia = &PLTError{message: "ERR_TABELA_VAZIA", status: http.StatusBadRequest}
var ErrDLNAAccess = &PLTError{message: "ERR_DLNA_NOT_RESPONDING", status: http.StatusInternalServerError}
var InvalidArtist = &PLTError{message: "ERR_INVALID_ARTIST", status: http.StatusBadRequest}
var InvalidGAA = &PLTError{message: "ERR_INVALID_GAA", status: http.StatusBadRequest}
var SongIsFavorite = &PLTError{message: "ERR_SONG_IS_FAVORITE", status: http.StatusConflict}
