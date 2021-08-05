#include <libdragon.h>
#include <stdbool.h>

void main(void)
{
	interrupts_init();
	debug_init();

	console_init();
	console_set_debug(true);

	printf("Hello world!\n");

	while(1) {}
}
