package c

import "fmt"

const vec = `
typedef struct %v {
	%v* data;
	int length;
	int capacity;
} %v;

%v new_%v() {
	%v v;
	v.length = 0;
	v.capacity = 4;
	v.data = malloc(sizeof(%v) * v.capacity);
	return v;
}

static void %v_resize(%v* this, int capacity) {
	%v* data = realloc(this->data, sizeof(%v) * capacity);
	this->data = data;
	this->capacity = capacity;
}

void %v_push(%v* this, %v v) {
	if (this->capacity == this->length) {
		%v_resize(this, this->capacity * 2);
	}
	this->data[this->length++] = v;
}

void %v_set(%v* this, int i, %v v) {
	if (i < 0) {
		i += this->length;
	}
	this->data[i] = v;
}

%v %v_get(%v* this, int i) {
	if (i < 0) {
		i += this->length;
	}
	return this->data[i];
}

%v %v_literal(%v* data, int n) {
	%v this = new_%v();
	for (int i = 0; i < n; i++) {
		%v_set(&this, i, data[i]);
	}
	return this;
}
`

func implementVector(name, ctype string) string {
	return fmt.Sprintf(vec, name, ctype, name, name, name, name, ctype, name, name, ctype, ctype, name, name, ctype, name, name, name, ctype, ctype, name, name, name, name, ctype, name, name, name)
}
