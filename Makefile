LIBSECRET_FLAGS := $(shell pkg-config --cflags --libs libsecret-1 glib-2.0)

dbg: dbg.cpp
	$(CXX) -O3 -o $@ $^ $(LIBSECRET_FLAGS)

cmd/key/command_enumer.go:: cmd/key/main.go
	cd cmd/key && enumer -type command -text