file = ik.csv

ycol = V(c,n)
x = time

all: downsample interp

downsample:
	./downsample.exe -xcol time -ycol="${ycol}" -sep tab -o ds.tsv ik.tsv

interp:
	./downsample.exe -xcol time -ycol="${ycol}" -sep tab -o intrp.tsv -interp ik.tsv