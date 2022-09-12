import requests
import sys
import time
import threading

SUBSTRING_SIZE = 21
BATCH_SIZE = 1000
START_INDEX = 1940374000
END_INDEX = -1
BUFFER_SIZE = 50
NAME = 'main'

keepAlive = True
digitBuffer = []
class ThreadWithReturnValue(threading.Thread):
    def __init__(self, group=None, target=None, name=None,
                 args=(), kwargs={}, Verbose=None):
        threading.Thread.__init__(self, group, target, name, args, kwargs)
        self._return = None
    def run(self):
        # print(type(self._target))
        if self._target is not None:
            self._return = self._target(*self._args,
                                                **self._kwargs)
    def join(self, *args):
        threading.Thread.join(self, *args)
        return self._return

def startThreads():
  digitsThread = threading.Thread(target=bufferizeDigits)
  digitsThread.start()
  return [digitsThread]

def bufferizeDigits():
  global digitBuffer
  global keepAlive
  digitIndex = START_INDEX
  while keepAlive:
    if len(digitBuffer) < BUFFER_SIZE:
      threads = []
      for i in range(BUFFER_SIZE - len(digitBuffer)):
        t = ThreadWithReturnValue(target=lambda : getPiDigits(digitIndex, BATCH_SIZE))
        t.start()
        threads.append(t)
        digitIndex += BATCH_SIZE
      
      for t in  threads:
        piDigits = t.join()
        digitBuffer.append(piDigits)
      
      threads = []

def getParam(paramName):
  params = sys.argv
  if paramName in params:
    for i in range(len(params)):
      if params[i] == paramName:
        return params[i+1]
  
  return None

def handleParams():
  argSize = getParam('-sz')
  if argSize != None:
    global SUBSTRING_SIZE
    SUBSTRING_SIZE = int(argSize)
  
  batchSize = getParam('-bs')
  if batchSize != None:
    global BATCH_SIZE
    BATCH_SIZE = int(batchSize)

  startIndex = getParam('-st')
  if startIndex != None:
    global START_INDEX
    START_INDEX = int(startIndex)

  endIndex = getParam('-end')
  if endIndex != None:
    global END_INDEX
    END_INDEX = int(endIndex)
  


def miller_rabin(n: int, allow_probable: bool = False) -> bool:
  """Created by Nathan Damon, @bizzfitch on github"""
  """Deterministic Miller-Rabin algorithm for primes ~< 3.32e24.

  Uses numerical analysis results to return whether or not the passed number
  is prime. If the passed number is above the upper limit, and
  allow_probable is True, then a return value of True indicates that n is
  probably prime. This test does not allow False negatives- a return value
  of False is ALWAYS composite.

  Parameters
  ----------
  n : int
    The integer to be tested. Since we usually care if a number is prime,
    n < 2 returns False instead of raising a ValueError.
  allow_probable: bool, default False
    Whether or not to test n above the upper bound of the deterministic test.

  Raises
  ------
  ValueError

  Reference
  ---------
  https://en.wikipedia.org/wiki/Miller%E2%80%93Rabin_primality_test
  """
  if n == 2:
    return True
  if not n % 2 or n < 2:
    return False
  if n > 5 and n % 10 not in (1, 3, 7, 9):  # can quickly check last digit
    return False
  if n > 3_317_044_064_679_887_385_961_981 and not allow_probable:
    raise ValueError(
      "Warning: upper bound of deterministic test is exceeded. "
      "Pass allow_probable=True to allow probabilistic test. "
      "A return value of True indicates a probable prime."
    )
  # array bounds provided by analysis
  bounds = [
    2_047,
    1_373_653,
    25_326_001,
    3_215_031_751,
    2_152_302_898_747,
    3_474_749_660_383,
    341_550_071_728_321,
    1,
    3_825_123_056_546_413_051,
    1,
    1,
    318_665_857_834_031_151_167_461,
    3_317_044_064_679_887_385_961_981,
  ]

  primes = [2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41]
  for idx, _p in enumerate(bounds, 1):
    if n < _p:
      # then we have our last prime to check
      plist = primes[:idx]
      break
  d, s = n - 1, 0
  # break up n -1 into a power of 2 (s) and
  # remaining odd component
  # essentially, solve for d * 2 ** s == n - 1
  while d % 2 == 0:
    d //= 2
    s += 1
  for prime in plist:
    pr = False
    for r in range(s):
      m = pow(prime, d * 2 ** r, n)
      # see article for analysis explanation for m
      if (r == 0 and m == 1) or ((m + 1) % n == 0):
        pr = True
        # this loop will not determine compositeness
        break
    if pr:
      continue
    # if pr is False, then the above loop never evaluated to true,
    # and the n MUST be composite
    return False
  return True

def getPiDigits(start, digitCount):
  # start is zero indexed
  params = {
    'start': start,
    'numberOfDigits': digitCount
  }
  pi_delivery_api = 'https://api.pi.delivery/v1/pi'
  shouldBreak = False
  while not shouldBreak:
    try:
      ans = requests.get(pi_delivery_api, params)
      jsonAns = ans.json()
      shouldBreak = True
      return jsonAns['content']
    except:
      print('error happened trying in 5 seconds')
      time.sleep(5)

def isPalindrome(s):
  isPalindrome = True
  
  for i in range(len(s)):
    if s[i] != s[len(s) - 1 - i]:
      isPalindrome = False
      break
  
  return isPalindrome

def iterativeCompare(initPattern, text):
  idx = 0
  pattern = initPattern
  while len(pattern) < SUBSTRING_SIZE:
    pattern += text[idx]
    idx += 1

  while idx < len(text):
    if isPalindrome(pattern): 
      if miller_rabin(int(pattern)):
        return True, pattern
    
    pattern = pattern[1:] + text[idx]
    idx += 1
  
  return False, pattern

def streamPiAfterPrimes():
  global digitBuffer
  global keepAlive
  begin = START_INDEX
  sz = BATCH_SIZE
  pattern = ""
  while True:
    with open('./parsed_index_' + NAME + '.txt', 'w')as f:
      f.flush()
      f.write('last index: ' + str(max(0, begin - sz)))
    if END_INDEX > 0 and begin > END_INDEX:
      keepAlive = False
      print('No value found until limit', END_INDEX)
      break
    
    print('verifying buffer for indexes', begin, ',', begin + sz - 1)
    while len(digitBuffer) == 0:
      pass
    print('trying indexes', begin, ',', begin + sz - 1)
    chunk = digitBuffer.pop(0)

    found, pattern = iterativeCompare(pattern, chunk)
    if found:
      keepAlive = False
      print("found the first polinomial prime number:")
      print(pattern)
      break
    begin += sz
  
  return pattern

if __name__ == "__main__":
  handleParams()
  threads = startThreads()
  streamPiAfterPrimes()
  pass

# last: 17073000