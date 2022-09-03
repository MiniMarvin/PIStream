# Palindromo em PI

Para resolver o problema do SigmaGeek a ideia é implementar algumas etapas
1. Encontrar todos os números primos de 9 digitos
2. Implementar Algoritmo para calcular vários digitos de PI
3. Algoritmo para consumir os digitos de PI e buscar uma sequência palindrômica sobre esses valores
4. Buscar os números primos na lista previamente gerada

Para determinar os números primos o que será determinado é o crivo de aristóteles.  

Para PI o algoritmo utilizado será a série de Chudnovsky dada por:
$$ \frac{1}{\pi} = 12 \sum^\infty_{k=0} \frac{(-1)^k (6k)! (13591409 + 545140134k)}{(3k)!(k!)^3 640320^{3k + 3/2}} $$
