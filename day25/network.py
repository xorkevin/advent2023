import networkx as nx

G = nx.Graph()
edges = []

with open("input.txt") as file:
    for line in file.readlines():
        l = line.rstrip().split(": ")
        lhs = l[0]
        rhs = l[1]
        for i in rhs.split(" "):
            G.add_edge(lhs, i, weight=1)
            edges.append([lhs, i])

cut_value, partition = nx.stoer_wagner(G)
print(cut_value, len(partition[0]) * len(partition[1]))
a = set(partition[0])
b = set(partition[1])
print([i for i in edges if (i[0] in a and i[1] in b) or (i[0] in b and i[1] in a)])
