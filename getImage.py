import math
import matplotlib.pyplot as plt
filename = "./att48.tsp"
edgeFile = open(filename)

xVals= []
yVals= []
maxX = -9999
minX = 9999

minY = 9999
maxY = -9999

solution0 = [6, 42, 29, 5, 35, 45, 32, 19, 46, 20, 12, 13, 24, 38, 31, 23, 9, 34, 44, 25, 3, 1, 28, 4, 41, 47, 33, 40, 15, 21, 2, 22, 10, 11, 14, 39, 8, 0, 7, 37, 30, 43, 17, 6, 27, 36, 18, 26]
solution1 = [0, 7, 37, 30, 43, 17, 6, 27, 5, 36, 18, 26, 16, 42, 29, 35, 45, 32, 19, 46, 20, 31, 38, 47, 4, 41, 23, 9, 44, 34, 3, 25, 1, 28, 33, 40, 15, 21, 2, 22, 13, 24, 12, 10, 11, 14, 39, 8]


for line in edgeFile:
    if not line.startswith("#"):
        words = line.split()
        x = float(words[1])
        y = float(words[2])

        if x < minX:
            minX = x
        elif x > maxX:
            maxX = x

        if y < minY:
            minY = y
        elif y > maxY:
            maxY = y

        xVals.append(x)
        yVals.append(y)

path0X = []
path0Y = []
lenPath0 = 0
for i in solution0:
    if len(path0X) > 0:
        x0 = path0X[-1]
        y0 = path0Y[-1]
        x1 = xVals[i]
        y1 = yVals[i]
        lenPath0 += math.hypot(x1 - x0, y1 - y1)
    path0X.append(xVals[i])
    path0Y.append(yVals[i])

path1X = []
path1Y = []
lenPath1 = 0
for i in solution1:
    if len(path0X) > 0:
        x0 = path0X[-1]
        y0 = path0Y[-1]
        x1 = xVals[i]
        y1 = yVals[i]
        lenPath1 += math.hypot(x1 - x0, y1 - y1)
    path1X.append(xVals[i])
    path1Y.append(yVals[i])


# uncomment to get numbered points

    #
    
f, (plt1, plt2)= plt.subplots(2, sharex=True)

for x, y, i in zip(xVals, yVals, range(0, 128)):
    plt1.text(x, y, i, color="blue", fontsize=10)

plt1.plot(path0X, path0Y, color = "blue", alpha = 0.5)
plt2.plot(path1X, path1Y, color = "red", alpha = 0.5)

plt1.plot(xVals, yVals, 'r.')
plt2.plot(xVals, yVals, 'b.')

plt1.set_title("Gentior Solution\nLength: " + str(lenPath0))
plt2.set_title("Known Optimal Solution\nLength: " + str(lenPath1))

#plt.axis([minX, maxX, minY, maxY])
plt.show()
