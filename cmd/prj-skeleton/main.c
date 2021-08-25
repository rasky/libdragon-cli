#include <libdragon.h>
#include <stdbool.h>
#include <stdio.h>

int main(void)
{
	init_interrupts();

	debug_init_usblog();   // debug console via USB (64drive / Everdrive)
	debug_init_isviewer(); // debug console on emulators

	console_init();
	console_set_debug(true);

	// Dump the contents of the file to the screen (on stdout)
	char buf[4096] = {0};

	FILE *f = fopen("rom:/credits.txt", "rb");
	fread(buf, 1, sizeof(buf), f);
	fclose(f);

	printf("%s", buf);

	while(1) {}
}
