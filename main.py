import sys
import traceback
from .pichudnovsky import PiChudnovsky

if __name__ == '__main__':
    try:
        if len(sys.argv) < 2:
            digits = 100
        else:
            digits = int(sys.argv[1])
        print("#### PI COMPUTATION ( {} digits )".format(digits))
        obj = PiChudnovsky(digits)
        tm = obj.compute()
        print("  Output  file:", "pi.txt")
        print("  Elapsed time: {} seconds".format(tm))
    except Exception as e:
        traceback.print_exc()
        sys.exit(1)