import requests
import math


def isPalindrome(s):
  isPalindrome = True
  
  for i in range(len(s)):
    if s[i] != s[len(s) - 1 - i]:
      isPalindrome = False
      break
  
  return isPalindrome

def generatePalindromesWithNineNumbers():
  palindromes = []
  for n in range(1, 99999 + 1):  
    s = str(n)
    if len(s) < 5:
      s += "".join(["0" for i in range(5 - len(s))])
    s += "".join(s[3::-1])
    palindromes.append(s)
  return [int(s) for s in palindromes]

def generateNineDigitsPalindromicPrimes():
  top = math.ceil(math.sqrt(10**10))
  palindromes = generatePalindromesWithNineNumbers()
  
  nums = [n for n in range(2, top + 1)]
  while nums:
    prime = nums[0]
    tmpNums = []
    for i in range(1, len(nums)):
      if nums[i]%prime != 0:
        tmpNums.append(nums[i])
    nums = tmpNums

    tmpPalindromes = []
    for i in range(0, len(palindromes)):
      if palindromes[i]%prime != 0:
        tmpPalindromes.append(palindromes[i])
    palindromes = tmpPalindromes
  
  return set([str(n) for n in palindromes])

def iterativeCompare(initPattern, text, searchSet):
  idx = 0
  pattern = initPattern
  while len(pattern) < 9:
    if str(chr(text[idx])) in [str(i) for i in range(10)]:
      pattern += str(chr(text[idx]))
    idx += 1

  while idx < len(text):
    # print(pattern, searchSet)
    if pattern in searchSet:
      return True, pattern
    
    if str(chr(text[idx])) in [str(i) for i in range(10)]:
      pattern = pattern[1:] + chr(text[idx])
    idx += 1
  # exit()
  
  return False, pattern

def streamPiAfterPrimes(searchSet):
  piUrl = "https://stuff.mit.edu/afs/sipb/contrib/pi/pi-billion.txt"

  s = requests.Session()
  r = s.get(piUrl, stream=True, verify=False)

  pattern = ""
  begin = 0
  for chunk in r.iter_content(1000):
    print('trying indexes', begin, ',', begin + 1000 - 1)
    begin += 1000
    found, pattern = iterativeCompare(pattern, chunk, searchSet)
    if found:
      print("found the first polinomial prime number:")
      print(pattern)
      break
  
  return pattern

if __name__ == "__main__":
  print('indexing primes...')
  primeSet = generateNineDigitsPalindromicPrimes()
  print('starting straming of text of pi')
  streamPiAfterPrimes(primeSet)
