1. C program that takes a number as input: 
```
#include <stdio.h>

int main() {
    int number;
    printf("Enter a number: ");
    scanf("%d", &number);
    printf("You entered: %d\n", number);

    return 0;
}

> Enter a number: You entered: XX
```

2. C program that runs indefinitely:
```
#include <stdio.h>

int main() {
    while(1) {
        printf("This program runs forever.\n");
    }
}
```