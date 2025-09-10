#include <stdlib.h>
#include <stdarg.h>
#include <stdio.h>
#include <string.h>

// example:
// llgoLoadPyModSyms(mod, "name1", &func1, "name2.attr.__get__", &func2, NULL)

typedef struct PyObject PyObject;

PyObject* PyObject_GetAttrString(PyObject* mod, const char* attrName);
void Py_IncRef(PyObject* o);
void Py_DecRef(PyObject* o);

PyObject* get_nested_attribute(PyObject* mod, const char* path) {
    PyObject* current = mod;
    Py_IncRef(current);
    const char* start = path;
    const char* end;

    while ((end = strchr(start, '.')) != NULL) {
        int size = end - start;
        char attr_name[size + 1];
        memcpy(attr_name, start, size);
        attr_name[size] = '\0';

        PyObject* next = PyObject_GetAttrString(current, attr_name);
        Py_DecRef(current);
        if (next == NULL) {
            return NULL;
        }

        current = next;
        start = end + 1;
    }

    PyObject* result = PyObject_GetAttrString(current, start);
    Py_DecRef(current);
    return result;
}

void llgoLoadPyModSyms(PyObject* mod, ...) {
    va_list ap;
    va_start(ap, mod);
    for (;;) {
        const char* name = va_arg(ap, const char*);
        if (name == NULL) {
            break;
        }
        PyObject** pfunc = va_arg(ap, PyObject**);
        if (*pfunc == NULL) {
            *pfunc = get_nested_attribute(mod, name);
        }
    }
    va_end(ap);
}

/*
wchar_t* toWcs(const char* str) {
    size_t len = mbstowcs(NULL, str, 0);
    wchar_t* wstr = (wchar_t*)malloc((len + 1) * sizeof(wchar_t));
    mbstowcs(wstr, str, len + 1);
    return wstr;
}

char* toMbs(const wchar_t* str) {
    size_t len = wcstombs(NULL, str, 0);
    char* mstr = (char*)malloc(len + 1);
    wcstombs(mstr, str, len + 1);
    return mstr;
}

wchar_t *Py_GetPath();

void Py_SetPath(const wchar_t* path);
void Py_Initialize();

void llgoPyInitialize() {
    setenv("PYTHONPATH", "/opt/homebrew/lib/python3.12/site-packages", 1);
    Py_Initialize();
    printf("sys.path = %s\n", toMbs(Py_GetPath()));
}
*/
