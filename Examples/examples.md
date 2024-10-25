# Code Snippet Examples
---

## Simple Java Program:
```
public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
    }
}

```

## Integer Input needed Java Program:
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
---

## Simple JavaScript Program:
```
function sayHello() {
    console.log("Hello, World!");
}

// Call the function
sayHello();

```

## Integer Input needed JavaScript Program:
```
const readline = require('readline');

// Create an interface for reading input from the user
const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

// Ask the user for a number
rl.question("Enter a number: ", function(input) {
    const number = parseInt(input, 10);

    // Check if the input is a valid number
    if (!isNaN(number)) {
        const result = number + 5;
        console.log("Result: " + result);
    } else {
        console.log("Please enter a valid number.");
    }

    // Close the readline interface
    rl.close();
});

```
---

## Simple Python Program:
```
def say_hello():
    print("Hello, World!")

# Call the function
say_hello()

```

## Integer Input needed Python Program:
```
def add_five(number):
    return number + 5

# Get user input
number = int(input("Enter a number: "))

# Call the function and print the result
print("Result:", add_five(number))

```
---

## Simple C Program:
```
#include <stdio.h>

int main() {
    printf("Hello, World!\n");
    return 0;
}

```

## Integer Input needed C Program:
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
---

## Simple C++ Program:
```
#include <iostream>

int main() {
    std::cout << "Hello, World!" << std::endl;
    return 0;
}

```

## Integer Input needed C++ Program:
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
---

6. Bash:

## Simple Bash Program:
```
#!/bin/bash
echo "Hello, World!"

```

## Integer Input needed Bash Program:
```
#!/bin/bash
read -p "Enter a number: " number
result=$((number + 5))
echo "Result: $result"

```
---

## C Program with Infinite Loop:
```
#include <stdio.h>

int main() {
    while(1) {
        printf("This program runs forever.\n");
    }
}
```