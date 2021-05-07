// See swig.org for more inteface options,
// e.g. map std::string to Go string
%module(directors="1") rtmlib
%feature("director") IRtmServiceEventHandler;
%feature("director") IWrappedIChannelEventHandler;
%{
#include "IAgoraRtmService.h"
#include "IAgoraRtmCallManager.h"
%}

%include "typemaps.i"
%include "std_string.i"
#ifdef SWIGGO
%apply (char* STRING, size_t LENGTH) { (const uint8_t* rawData, int length) }
#endif

%include "IAgoraRtmService.h"
%include "IAgoraRtmCallManager.h"