# Code Snippet Examples

1. Java:

## Simple
```
public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
    }
}

```

## Integer Input needed
```
import java.util.Scanner;

public class Main {
    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);
        System.out.print("Enter a number: ");
        int number = scanner.nextInt();
        int result = number + 5;
        System.out.println("Result: " + result);
    }
}
```

2. JavaScript:

## Simple
```
function sayHello() {
    console.log("Hello, World!");
}

// Call the function
sayHello();

```

## Integer Input needed
```
function addFive(number) {
    return number + 5;
}

// Prompt user for input
const number = parseInt(prompt("Enter a number: "), 10);

// Call the function and log the result
console.log("Result: " + addFive(number));

```

3. Python:

## Simple
```
def say_hello():
    print("Hello, World!")

# Call the function
say_hello()

```

## Integer Input needed
```
def add_five(number):
    return number + 5

# Get user input
number = int(input("Enter a number: "))

# Call the function and print the result
print("Result:", add_five(number))

```

4. C:

## Simple
```
#include <stdio.h>

int main() {
    printf("Hello, World!\n");
    return 0;
}

```

## Integer Input needed
```
#include <stdio.h>

int main() {
    int number;
    printf("Enter a number: ");
    scanf("%d", &number);
    int result = number + 5;
    printf("Result: %d\n", result);
    return 0;
}

```

5. C++:

## Simple
```
#include <iostream>

int main() {
    std::cout << "Hello, World!" << std::endl;
    return 0;
}

```

## Integer Input needed
```
#include <iostream>

int main() {
    int number;
    std::cout << "Enter a number: ";
    std::cin >> number;
    int result = number + 5;
    std::cout << "Result: " << result << std::endl;
    return 0;
}
```

6. Bash:

## Simple
```
#!/bin/bash
echo "Hello, World!"

```

## Integer Input needed
```
#!/bin/bash
read -p "Enter a number: " number
result=$((number + 5))
echo "Result: $result"

```
======================

## C Program with Infinite Loop:
```
#include <stdio.h>

int main() {
    while(1) {
        printf("This program runs forever.\n");
    }
}
```