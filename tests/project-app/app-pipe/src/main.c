#ifndef BUFFER_SIZE
#error "BUFFER_SIZE must be defined"
#endif

#include <stdio.h>

static constexpr int buffer_size = BUFFER_SIZE;
static_assert(buffer_size == 80, "Expected BUFFER_SIZE == 80");

int main(int argc, char *argv[]) {
    FILE *pipe_fp, *infile;
    static char readbuf[buffer_size];

    if (argc != 3) {
        fprintf(stderr, "USAGE: app-pipe [command] [filename]\n");
        return 1;
    }

    if ((infile = fopen(argv[2], "rt")) == NULL) {
        perror("fopen");
        return 1;
    }
    if ((pipe_fp = popen(argv[1], "w")) == NULL) {
        perror("popen");
        return 1;
    }
    do {
        fgets(readbuf, buffer_size, infile);
        if (feof(infile)) { break; }
        fputs(readbuf, pipe_fp);
    } while (!feof(infile));

    fclose(infile);
    pclose(pipe_fp);

    return 0;
}
