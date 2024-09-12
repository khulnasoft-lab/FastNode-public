from fastnode import Fastnode

class Man(object):
    def __init__(self):
        self.code = "man"
    def car(self):
        print("star!")

k = Fastnode()
k.foo()

m = Man()
m.car()

def print_code():
    print(k.code + " " + m.code)

def some_fastnode():
    return Fastnode()

print_code()

q = some_fastnode()