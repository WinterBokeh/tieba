CREATE TABLE userinfo(
    uid INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(10) NOT NULL UNIQUE,
    pwd VARCHAR(20) NOT NULL,
    email VARCHAR(20) NOT NULL,
    statement  VARCHAR(30),
    reg_data DATE NOT NULL,
    state INT NOT NULL DEFAULT 0,
    power INT NOT NULL DEFAULT 0
);

#include <iostream>
using namespace std;

int Sum(int n) {
	int sum = 0;
	for (int i = 1; i <= n; ++i) {
		sum += i;
	}
	return sum;
}

int Sub(int n) {
	int sum = 1;
	for (int i = 1; i <= n; ++i) {
		sum *= i;
	}
	return sum;
}

int main() {
	int n;
	cin >> n;
	cout << Sum(n) << endl;
	cout << Sub(n);
}