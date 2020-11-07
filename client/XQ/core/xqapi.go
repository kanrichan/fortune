package core

/*#include <stdio.h>
#include <string.h>
#include <windows.h>
#include <stdint.h>
#include <stdlib.h>

#define XQAPI(RetType, Name, ...)													\
	typedef RetType(__stdcall *Name##_Type)(unsigned char * authid, ##__VA_ARGS__); \
	Name##_Type Name##_Ptr;															\
	RetType Name(__VA_ARGS__);

#define LoadAPI(Name) Name##_Ptr = (Name##_Type)GetProcAddress(hmod, #Name)

unsigned char * authid;
XQAPI(void, S3_Api_SendMsg, char *, int, char *, char *, char *, int);
XQAPI(void, S3_Api_OutPutLog, char *);

extern void __stdcall XQ_AuthId(int ID, int IMAddr){
	authid = (unsigned char *)malloc(sizeof(unsigned char)*16);
	*((int*)authid) = 1;
	*((int*)(authid + 4)) = 8;
	*((int*)(authid + 8)) = ID;
	*((int*)(authid + 12)) = IMAddr;
	authid += 8;
	HMODULE hmod = LoadLibraryA("xqapi.dll");
	LoadAPI(S3_Api_SendMsg);
	LoadAPI(S3_Api_OutPutLog);
	return;
}

void S3_Api_SendMsg(char * var0, int var1, char * var2, char * var3, char * var4, int var5){
	S3_Api_SendMsg_Ptr(authid, var0, var1, var2, var3, var4, var5);
	free(var0);
	free(var2);
	free(var3);
	free(var4);
}

void S3_Api_OutPutLog(char * var0){
	S3_Api_OutPutLog_Ptr(authid, var0);
	free(var0);
}
*/
import "C"
import sc "golang.org/x/text/encoding/simplifiedchinese"

func CString(str string) *C.char {
	gbstr, _ := sc.GB18030.NewEncoder().String(str)
	return C.CString(gbstr)
}

func GoString(str *C.char) string {
	utf8str, _ := sc.GB18030.NewDecoder().String(C.GoString(str))
	return utf8str
}

func CBool(b bool) C.int32_t {
	if b {
		return 1
	}
	return 0
}
func SendPrivateMsg(var0 string, var1 int32, var2 string, var3 string, var4 string, var5 int32) {
	C.S3_Api_SendMsg(
		CString(var0), C.int(var1), CString(var2), CString(var3), CString(var4), C.int(var5),
	)
}

func OutPutLog(var0 string) {
	C.S3_Api_OutPutLog(
		CString(var0),
	)
}
