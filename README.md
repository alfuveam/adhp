# ADHP: AMBIENTE PARA O DESENVOLVIMENTO DE HABILIDADES EM PROGRAMAÇÃO

Bem Vindo!
Essa ferramenta tem a função de ajudar o instrutor, retirando o trabalho repetitivo de corrigir questões de programação de alunos que estão iniciando na área de programação. Além de ajudar os alunos com a utilização de feedback automatizado para realizar o termino dos exercicios de forma mais eficaz, sendo então explorado o uso da repetição espaçada.

# Rotas
    Docente Rota
<img src="./assets/docente_rota.png">

--------------------                              --------------------
    Discente Rota
<img src="./assets/discente_rota_no_submit_repeticao_espacada.png">

--------------------                              --------------------
    Discente Rota
<img src="./assets/discente_rota_lista_exercicio.png">

# Comandos

Para build

	docker compose -f docker-compose.dev.yml build

Para execução

	docker compose -f docker-compose.dev.yml up -d

# Go

module tcc_ead/sum go 1.23.0
```
// go.mod
module tcc_ead/sum

go 1.23.0
```
    Exercício exemplo Golang: 1

```
// sum.go
// this code is running correctly

package main

import "fmt"

func Sum(numOne int, numTwo int) int {
	if numOne == 0 && numTwo == 0 {
		return -1
	}

	if numOne < 0 || numTwo < 0 {
		return 0
	}
	return numOne + numTwo
}

// Você deve fazer uma função que some somente numeros positivos, se por algum motivo o usuario colocar
// 0 + 0, você deve retornar -1
func main() {
    fmt.Println(Sum(0, 0))
}
```

```
// sum_test.go
package main

import "testing"

func TestSumPositiveNumber(t *testing.T) {
	value := Sum(0, 1)

	if value != 1 {
		t.Fatalf("Logic to positive value it's wrong")
	}

	value = Sum(4, 5)

	if value != 9 {
		t.Fatalf("Logic to positive value it's wrong")
	}

	value = Sum(10, 5)

	if value != 15 {
		t.Fatalf("Logic to positive value it's wrong")
	}
}

func TestSumZero(t *testing.T) {
	value := Sum(0, 0)

	if value != -1 {
		t.Fatalf("Check description of exercise: %d", value)
	}
}

func TestSumNegativeNumber(t *testing.T) {
	value := Sum(0, -1)

	if value < 0 {
		t.Fatalf("Logic to negative value it's wrong")
	}

	value = Sum(-1, -1)

	if value < 0 {
		t.Fatalf("Logic to negative value it's wrong")
	}
}
```
--------------------                              --------------------
    Exercício exemplo Golang: 2

```
package main

import "fmt"

func main() {
	a := 5
	b := 5
	c := a + b
	fmt.Println(c)
}
```
```
package main

import (
	"bytes"
	"os"
	"testing"
)

// Função auxiliar para capturar a saída de uma função
func captureOutput(f func()) string {
	// Salva a saída padrão atual
	old := os.Stdout

	// Cria um novo buffer para capturar a saída
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Executa a função
	f()

	// Restaura a saída padrão
	w.Close()
	os.Stdout = old

	// Lê o conteúdo do buffer
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestMainOutput(t *testing.T) {
	// Captura a saída da função main
	output := captureOutput(main)

	// Remove a nova linha no final da saída
	output = output[:len(output)-1]

	// Verifica se a saída é "10"
	expected := "10"
	if output != expected {
		t.Errorf("Saída incorreta. Esperado: %s, Obtido: %s", expected, output)
	}
}
```
# Python
    Exercício exemplo Python: 1

```
# sum.py
def sum(num_one, num_two):
    if num_one == 0 and num_two == 0:
        return -1
    
    if num_one < 0 or num_two < 0:
        return 0
    
    return num_one + num_two

# Você deve fazer uma função que some somente numeros positivos, se por algum motivo o usuario colocar
# 0 + 0, você deve retornar -1
def main():
    result = sum(0, 0)
    print(result)  # This will print -1

if __name__ == "__main__":
    main()
```
```
# sum_test.py
import unittest

# The Sum function from the previous example
def sum(num_one, num_two):
    if num_one == 0 and num_two == 0:
        return -1
    
    if num_one < 0 or num_two < 0:
        return 0
    
    return num_one + num_two

# Test cases
class TestSum(unittest.TestCase):
    def test_sum_positive_numbers(self):
        self.assertEqual(sum(0, 1), 1, "Logic for positive values is wrong")
        self.assertEqual(sum(4, 5), 9, "Logic for positive values is wrong")
        self.assertEqual(sum(10, 5), 15, "Logic for positive values is wrong")

    def test_sum_zero(self):
        self.assertEqual(sum(0, 0), -1, "Check description of exercise")

    def test_sum_negative_numbers(self):
        self.assertEqual(sum(0, -1), 0, "Logic for negative values is wrong")
        self.assertEqual(sum(-1, -1), 0, "Logic for negative values is wrong")

# Run the tests
if __name__ == "__main__":
    unittest.main()
```

# Working
curl -X POST -H "Content-Type: application/json" -H "X-Auth-Token: S2AD13AS2D13AS21DA3SD1AS32D13SAD13A1SS5WE8EERRTYasdda6sd5467e7e7e7e" -d '
{
  "source_from_user": "ZGVmJTIwc3VtKG51bV9vbmUlMkMlMjBudW1fdHdvKSUzQSUwQSUwQSUyMGlmJTIwbnVtX29uZSUyMCUzRCUzRCUyMDAlMjBhbmQlMjBudW1fdHdvJTIwJTNEJTNEJTIwMCUzQSUwQSUwQSUyMHJldHVybiUyMC0xJTBBJTBBJTIwJTBBJTBBJTIwaWYlMjBudW1fb25lJTIwJTNDJTIwMCUyMG9yJTIwbnVtX3R3byUyMCUzQyUyMDAlM0ElMEElMEElMjByZXR1cm4lMjAwJTBBJTBBJTIwJTBBJTBBJTIwcmV0dXJuJTIwbnVtX29uZSUyMCUyQiUyMG51bV90d28lMEElMEElMEElMEElMEElMjMlMjBWb2MlQzMlQUElMjBkZXZlJTIwZmF6ZXIlMjB1bWElMjBmdW4lQzMlQTclQzMlQTNvJTIwcXVlJTIwc29tZSUyMHNvbWVudGUlMjBudW1lcm9zJTIwcG9zaXRpdm9zJTJDJTIwc2UlMjBwb3IlMjBhbGd1bSUyMG1vdGl2byUyMG8lMjB1c3VhcmlvJTIwY29sb2NhciUwQSUwQSUyMyUyMDAlMjAlMkIlMjAwJTJDJTIwdm9jJUMzJUFBJTIwZGV2ZSUyMHJldG9ybmFyJTIwLTElMEElMEFkZWYlMjBtYWluKCklM0ElMEElMEElMjByZXN1bHQlMjAlM0QlMjBzdW0oMCUyQyUyMDApJTBBJTBBJTIwcHJpbnQocmVzdWx0KSUyMCUyMyUyMFRoaXMlMjB3aWxsJTIwcHJpbnQlMjAtMSUwQSUwQSUwQSUwQSUwQWlmJTIwX19uYW1lX18lMjAlM0QlM0QlMjAlMjJfX21haW5fXyUyMiUzQSUwQSUwQSUyMG1haW4oKSUwQSUwQSUwQQ",
  "source_unit_teste": "aW1wb3J0JTIwdW5pdHRlc3QlMEElMEElMEElMEElMEElMjMlMjBUaGUlMjBTdW0lMjBmdW5jdGlvbiUyMGZyb20lMjB0aGUlMjBwcmV2aW91cyUyMGV4YW1wbGUlMEElMEFkZWYlMjBzdW0obnVtX29uZSUyQyUyMG51bV90d28pJTNBJTBBJTBBJTIwJTIwJTIwJTIwaWYlMjBudW1fb25lJTIwJTNEJTNEJTIwMCUyMGFuZCUyMG51bV90d28lMjAlM0QlM0QlMjAwJTNBJTBBJTBBJTIwJTIwJTIwJTIwJTIwJTIwJTIwJTIwcmV0dXJuJTIwLTElMEElMEElMjAlMjAlMjAlMjAlMEElMEElMjAlMjAlMjAlMjBpZiUyMG51bV9vbmUlMjAlM0MlMjAwJTIwb3IlMjBudW1fdHdvJTIwJTNDJTIwMCUzQSUwQSUwQSUyMCUyMCUyMCUyMCUyMCUyMCUyMCUyMHJldHVybiUyMDAlMEElMEElMjAlMjAlMjAlMjAlMEElMEElMjAlMjAlMjAlMjByZXR1cm4lMjBudW1fb25lJTIwJTJCJTIwbnVtX3R3byUwQSUwQSUwQSUwQSUwQSUyMyUyMFRlc3QlMjBjYXNlcyUwQSUwQWNsYXNzJTIwVGVzdFN1bSh1bml0dGVzdC5UZXN0Q2FzZSklM0ElMEElMEElMjAlMjAlMjAlMjBkZWYlMjB0ZXN0X3N1bV9wb3NpdGl2ZV9udW1iZXJzKHNlbGYpJTNBJTBBJTBBJTIwJTIwJTIwJTIwJTIwJTIwJTIwJTIwc2VsZi5hc3NlcnRFcXVhbChzdW0oMCUyQyUyMDEpJTJDJTIwMSUyQyUyMCUyMkxvZ2ljJTIwZm9yJTIwcG9zaXRpdmUlMjB2YWx1ZXMlMjBpcyUyMHdyb25nJTIyKSUwQSUwQSUyMCUyMCUyMCUyMCUyMCUyMCUyMCUyMHNlbGYuYXNzZXJ0RXF1YWwoc3VtKDQlMkMlMjA1KSUyQyUyMDklMkMlMjAlMjJMb2dpYyUyMGZvciUyMHBvc2l0aXZlJTIwdmFsdWVzJTIwaXMlMjB3cm9uZyUyMiklMEElMEElMjAlMjAlMjAlMjAlMjAlMjAlMjAlMjBzZWxmLmFzc2VydEVxdWFsKHN1bSgxMCUyQyUyMDUpJTJDJTIwMTUlMkMlMjAlMjJMb2dpYyUyMGZvciUyMHBvc2l0aXZlJTIwdmFsdWVzJTIwaXMlMjB3cm9uZyUyMiklMEElMEElMEElMEElMEElMjAlMjAlMjAlMjBkZWYlMjB0ZXN0X3N1bV96ZXJvKHNlbGYpJTNBJTBBJTBBJTIwJTIwJTIwJTIwJTIwJTIwJTIwJTIwc2VsZi5hc3NlcnRFcXVhbChzdW0oMCUyQyUyMDApJTJDJTIwLTElMkMlMjAlMjJDaGVjayUyMGRlc2NyaXB0aW9uJTIwb2YlMjBleGVyY2lzZSUyMiklMEElMEElMEElMEElMEElMjAlMjAlMjAlMjBkZWYlMjB0ZXN0X3N1bV9uZWdhdGl2ZV9udW1iZXJzKHNlbGYpJTNBJTBBJTBBJTIwJTIwJTIwJTIwJTIwJTIwJTIwJTIwc2VsZi5hc3NlcnRFcXVhbChzdW0oMCUyQyUyMC0xKSUyQyUyMDAlMkMlMjAlMjJMb2dpYyUyMGZvciUyMG5lZ2F0aXZlJTIwdmFsdWVzJTIwaXMlMjB3cm9uZyUyMiklMEElMEElMjAlMjAlMjAlMjAlMjAlMjAlMjAlMjBzZWxmLmFzc2VydEVxdWFsKHN1bSgtMSUyQyUyMC0xKSUyQyUyMDAlMkMlMjAlMjJMb2dpYyUyMGZvciUyMG5lZ2F0aXZlJTIwdmFsdWVzJTIwaXMlMjB3cm9uZyUyMiklMEElMEElMEElMEElMEElMjMlMjBSdW4lMjB0aGUlMjB0ZXN0cyUwQSUwQWlmJTIwX19uYW1lX18lMjAlM0QlM0QlMjAlMjJfX21haW5fXyUyMiUzQSUwQSUwQSUyMCUyMCUyMCUyMHVuaXR0ZXN0Lm1haW4oKQ=",
  "lista": "1",
  "exercicio": "1",
  "usuario": "123456789"
}
' localhost:8082/api/v1/run-test-python

# Working
curl -X POST -H "Content-Type: application/json" -d '
{
  "source_from_user": "Ly8gc3VtLmdvDQpwYWNrYWdlIG1haW4NCg0KZnVuYyBTdW0obnVtT25lIGludCwgbnVtVHdvIGludCkgaW50IHsNCglpZiBudW1PbmUgPT0gMCAmJiBudW1Ud28gPT0gMCB7DQoJCXJldHVybiAtMQ0KCX0NCg0KCWlmIG51bU9uZSA8IDAgfHwgbnVtVHdvIDwgMCB7DQoJCXJldHVybiAwDQoJfQ0KCXJldHVybiBudW1PbmUgKyBudW1Ud28NCn0NCg0KLy8gVm9jw6ogZGV2ZSBmYXplciB1bWEgZnVuw6fDo28gcXVlIHNvbWUgc29tZW50ZSBudW1lcm9zIHBvc2l0aXZvcywgc2UgcG9yIGFsZ3VtIG1vdGl2byBvIHVzdWFyaW8gY29sb2Nhcg0KLy8gMCArIDAsIHZvY8OqIGRldmUgcmV0b3JuYXIgLTENCmZ1bmMgbWFpbigpIHsNCglfID0gU3VtKDAsIDApDQp9",
  "source_unit_teste": "cGFja2FnZSBtYWluDQoNCmltcG9ydCAidGVzdGluZyINCg0KZnVuYyBUZXN0U3VtUG9zaXRpdmVOdW1iZXIodCAqdGVzdGluZy5UKSB7DQoJdmFsdWUgOj0gU3VtKDAsIDEpDQoNCglpZiB2YWx1ZSAhPSAxIHsNCgkJdC5GYXRhbGYoIkxvZ2ljIHRvIHBvc2l0aXZlIHZhbHVlIGl0J3Mgd3JvbmciKQ0KCX0NCg0KCXZhbHVlID0gU3VtKDQsIDUpDQoNCglpZiB2YWx1ZSAhPSA5IHsNCgkJdC5GYXRhbGYoIkxvZ2ljIHRvIHBvc2l0aXZlIHZhbHVlIGl0J3Mgd3JvbmciKQ0KCX0NCg0KCXZhbHVlID0gU3VtKDEwLCA1KQ0KDQoJaWYgdmFsdWUgIT0gMTUgew0KCQl0LkZhdGFsZigiTG9naWMgdG8gcG9zaXRpdmUgdmFsdWUgaXQncyB3cm9uZyIpDQoJfQ0KfQ0KDQpmdW5jIFRlc3RTdW1aZXJvKHQgKnRlc3RpbmcuVCkgew0KCXZhbHVlIDo9IFN1bSgwLCAwKQ0KDQoJaWYgdmFsdWUgIT0gLTEgew0KCQl0LkZhdGFsZigiQ2hlY2sgZGVzY3JpcHRpb24gb2YgZXhlcmNpc2U6ICVkIiwgdmFsdWUpDQoJfQ0KfQ0KDQpmdW5jIFRlc3RTdW1OZWdhdGl2ZU51bWJlcih0ICp0ZXN0aW5nLlQpIHsNCgl2YWx1ZSA6PSBTdW0oMCwgLTEpDQoNCglpZiB2YWx1ZSA8IDAgew0KCQl0LkZhdGFsZigiTG9naWMgdG8gbmVnYXRpdmUgdmFsdWUgaXQncyB3cm9uZyIpDQoJfQ0KDQoJdmFsdWUgPSBTdW0oLTEsIC0xKQ0KDQoJaWYgdmFsdWUgPCAwIHsNCgkJdC5GYXRhbGYoIkxvZ2ljIHRvIG5lZ2F0aXZlIHZhbHVlIGl0J3Mgd3JvbmciKQ0KCX0NCn0=",
  "lista": "1",
  "exercicio": "1",
  "usuario": "123456789"
}
' localhost:8082/api/v1/run-test-golang
