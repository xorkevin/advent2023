import z3

lines = []
with open("input.txt", "r") as file:
    for i in range(3):
        line = file.readline().rstrip().split(" @ ")
        lines.append(
            [
                [int(i) for i in line[0].split(", ")],
                [int(i) for i in line[1].split(", ")],
            ]
        )

xarr = [z3.Real("x{}".format(i)) for i in range(3)]
varr = [z3.Real("v{}".format(i)) for i in range(3)]
tarr = [z3.Real("t{}".format(i)) for i in range(3)]

s = z3.Solver()

for n, line in enumerate(lines):
    t = tarr[n]
    for i in range(3):
        s.add(line[0][i] + line[1][i] * t == xarr[i] + varr[i] * t)


print(s.check())
model = s.model()
print(model.eval(xarr[0] + xarr[1] + xarr[2]))
