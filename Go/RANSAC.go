// Name: Anouar Smaili

package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"math/rand"
	"sync"
	"time"
	"path/filepath"
    "strings"
	//"runtime"
)

type Point3D struct {
	X float64
	Y float64
	Z float64
}

type Plane3D struct {
	A float64
	B float64
	C float64
	D float64
}

type Plane3DwSupport struct {
	Plane3D
	SupportSize int
}

// reads an XYZ file and returns a slice of Point3D
func ReadXYZ(filename string) []Point3D {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	var points []Point3D
	var x, y, z float64
	var xSkip, ySkip, zSkip string
	// Skip the line "x y z"
	fmt.Fscanf(file, "%s %s %s\n", &xSkip, &ySkip, &zSkip)
	for {
		_, err := fmt.Fscanf(file, "%f %f %f\n", &x, &y, &z)
		if err != nil {
			break
		}
		points = append(points, Point3D{X: x, Y: y, Z: z})
	}

	return points
}

// saves a slice of Point3D into an XYZ file
func SaveXYZ(filename string, points []Point3D) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
	defer file.Close()

	for _, p := range points {
		fmt.Fprintf(file, "%f %f %f\n", p.X, p.Y, p.Z)
	}
}

// computes the distance between points p1 and p2
func (p1 *Point3D) GetDistance(p2 *Point3D) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	dz := p1.Z - p2.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// computes the plane defined by a slice of 3 points
func GetPlane(points []Point3D) Plane3D {

	x1 := points[0].X
	y1 := points[0].Y
	z1 := points[0].Z
	x2 := points[1].X
	y2 := points[1].Y
	z2 := points[1].Z
	x3 := points[2].X
	y3 := points[2].Y
	z3 := points[2].Z

    a1 := x2 - x1
    b1 := y2 - y1
    c1 := z2 - z1
    a2 := x3 - x1
    b2 := y3 - y1
    c2 := z3 - z1
    a := b1 * c2 - b2 * c1
    b := a2 * c1 - a1 * c2
    c := a1 * b2 - b1 * a2
    d := (- a * x1 - b * y1 - c * z1)

	return Plane3D{a, b, c, d}
}

// computes the number of required RANSAC iterations
func GetNumberOfIterations(confidence float64, percentageOfPointsOnPlane float64) int{
	return int(math.Log(1 - confidence) / math.Log(1-math.Pow(percentageOfPointsOnPlane, 3)))
}

// Returns the distance from a point to the plane
func (p Plane3D) GetDistanceFromPlane(pt Point3D) float64 {
    m := math.Abs(p.A*pt.X + p.B*pt.Y + p.C*pt.Z + p.D)
    n := math.Sqrt(p.A*p.A + p.B*p.B + p.C*p.C)
    return m / n
}

// computes the support of a plane in a slice of points
func GetSupport(plane Plane3D, points []Point3D, eps float64) Plane3DwSupport{
	SupportSize := 0
	for _,point := range points {
		// Check if the point belongs to the plane
		if (plane.GetDistanceFromPlane(point)) <= eps {
			SupportSize ++
		}
	}
	return Plane3DwSupport{plane, SupportSize}
}

// extracts the points that supports the given plane
// and returns them in a slice of points
func GetSupportingPoints(plane Plane3D, points []Point3D, eps float64) []Point3D {
	var supportPoints []Point3D
	for _,point := range points {
		// Check if the point belongs to the plane
		if (plane.GetDistanceFromPlane(point)) <= eps {
			supportPoints = append(supportPoints, point)
		}
	}
	return supportPoints
}

// creates a new slice of points in which all points
// belonging to the plane have been removed
func RemovePlane(plane Plane3D, points []Point3D, eps float64) []Point3D {
	var newPoints []Point3D
	for _,point := range points {
		// Check if the point belongs to the plane
		if !((plane.GetDistanceFromPlane(point)) <= eps) {
			newPoints = append(newPoints, point)
		}
	}
	return newPoints
}

// Generate random number from 0 to size-1 (inclus)
func randomNumbers(size int) int {
        
	return rand.Intn(size)
}

// Random point generator : -> Point3D 
// It randomly selects a point from the provided slice of Point3D (the input point cloud). 
// Its output channel transmits instances of Point3D. 
func randomGenerator(wg *sync.WaitGroup, stop <-chan bool, fct func(int) int, points []Point3D) <-chan Point3D {

	intStream := make(chan Point3D)
	size := len(points)
	
	go func() {
		defer func() {wg.Done()}()
		defer close(intStream)
		
		for {
		   select {
		      case <-stop:
		         return
		      case intStream <- points[fct(size)]:
		   }
		}
	}()
	
	return intStream
}

// Triplet of points generator: Point3D -> [3]Point3D 
// It reads Point3D instances from its input channel and accumulate 3 points. 
// Its output channel transmits arrays of Point3D (composed of three points). 
func triplet(wg *sync.WaitGroup, stop <-chan bool, inputIntstream <-chan Point3D, iterations int) <-chan [3]Point3D {

	outputIntStream := make(chan [3]Point3D)
	
	go func() {
		defer func() {wg.Done()}()
		defer close(outputIntStream)

		count := 0
		var ok bool
		var point Point3D
		var triplet [3]Point3D
		// Loop multiple of 3 times depending on the number of iterations
		for count < (iterations*3){
		   	select {
		    	case <-stop:
		        	break
				case point, ok= <- inputIntstream:
					if !ok {
						return
					}
					triplet[count%3]= point
					count ++
					if count%3==0 {
						outputIntStream <- triplet
						triplet = [3]Point3D{}
					}
		   }
		}
	}()
	
	return outputIntStream
}

// TakeN: [3]Point3D -> [3]Point3D 
// It reads arrays of Point3D and resend them. 
// It automatically stops the pipeline after having received N arrays. 
func takeN(wg *sync.WaitGroup, stop <-chan bool, inputIntstream <-chan [3]Point3D, n int) <-chan [3]Point3D {

	outputIntStream := make(chan [3]Point3D)
	
	go func() {
		defer func() {wg.Done()}()
		defer close(outputIntStream)
		for i:=0; i<n; i++ {
			select {
				case <-stop:
					break
				case outputIntStream <- <- inputIntstream:
			}
		}
	}()
	
	return outputIntStream
}

// Plane estimator: [3]Point3D -> Plane3D 
// It reads arrays of three Point3D and compute the plane defined by these points. 
// Its output channel transmits Plane3D instances describing the computed plane parameters
func planeEstimator(wg *sync.WaitGroup, stop <-chan bool, inputIntstream <-chan [3]Point3D, fct func([]Point3D) Plane3D) <-chan Plane3D {

	outputIntStream := make(chan Plane3D)
	
	go func() {
		defer func() {wg.Done()}()
		defer close(outputIntStream)
		var table [3]Point3D
		var ok bool
		for {
			select {
				case <-stop:
					return
				case table, ok= <- inputIntstream:
					if !ok {
						return
					}
					slice := []Point3D{}
					slice = append(slice, table[0], table[1], table[2])
					outputIntStream <- fct(slice)
			}
		}
	}()
	
	return outputIntStream
}

// Supporting point finder: Plane3D -> Plane3DwSupport 
// It counts the number of points in the provided slice of Point3D (the input point cloud) that supports the receive 3D plane. 
// Its output channel transmits the plane parameters and the number of supporting points in a Point3DwSupport instance.  
func supportingPointFinder(wg *sync.WaitGroup, stop <-chan bool, inputIntstream <-chan Plane3D, fct func(Plane3D, []Point3D, float64) Plane3DwSupport, eps float64, originalPoints []Point3D) <-chan Plane3DwSupport {

	outputIntStream := make(chan Plane3DwSupport)
	
	go func() {
		defer func() {wg.Done()}()
		defer close(outputIntStream)
		var plane Plane3D
		var ok bool
		for {
			select {
				case <-stop:
					return
				case plane, ok= <- inputIntstream:
					if !ok {
						return
					}
					outputIntStream <- fct(plane, originalPoints, eps)
			}
		}
	}()
	
	return outputIntStream
}

// Fan In: Plane3DwSupport -> Plane3DwSupport 
// It multiplexes the results receives from multiple channels into one output channel. 
// Cette composante multiplexe les résultats recu de multiples channels dans un seul channel de sortie.
func fanIn(wg *sync.WaitGroup, stop <-chan bool, channels []<-chan Plane3DwSupport) <-chan Plane3DwSupport {

    var multiplexGroup sync.WaitGroup
	outputIntStream := make(chan Plane3DwSupport)
	
	reader:= func(ch <-chan Plane3DwSupport) {
		defer func() {multiplexGroup.Done()}()
		for i:= range ch {
		   
		   select {
		      case <-stop:
		         return
		      case outputIntStream <- i:
		   }
		}
	}
	
	// all goroutine must return before 
	// the output channel is closed
	multiplexGroup.Add(len(channels))
	for _, ch := range channels {
	
	   go reader(ch)
	}
	
	go func() {
	
	   defer func() {wg.Done()}()
	   defer close(outputIntStream)
	   multiplexGroup.Wait()
	}()

	return outputIntStream
}

// Dominant plane identifier: Plane3DwSupport 
// It receives Plane3DwSupport instances and keep in memory the plane with the best support received so far. 
// This component does not output values, it simply maintains the provided *Plane3DwSupport variable.
func DominantPlanIdentifier(wg *sync.WaitGroup, stop <-chan bool, inputIntstream <-chan Plane3DwSupport, bestCurrentSupport *Plane3DwSupport) <-chan int{
	
	outputIntStream := make(chan int)

	go func() {
		defer func() {wg.Done()}()
		defer close(outputIntStream)
		var support Plane3DwSupport
		var ok bool
		for {
			select {
				case <-stop:
					return
				case support, ok= <- inputIntstream:
					if !ok {
						return
					}
					if (support.SupportSize > bestCurrentSupport.SupportSize){
						bestCurrentSupport.Plane3D = support.Plane3D
						bestCurrentSupport.SupportSize = support.SupportSize
					}
			}
		}
	}()
	return outputIntStream
}


func main() {
	// Read the XYZ file specified as a first argument to your go program and create the corresponding
	// slice of Point3D composed of the set of points from the XYZ file.
	var originalPoints []Point3D
	originalPoints = ReadXYZ(os.Args[1])

	// Create a bestSupport variable of type Plane3DwSupport initialized to all 0s.
	var bestSupport Plane3DwSupport

	// Find the number of iterations required based on the specified confidence and percentage
	// provided as 2nd and 3rd argument.
	confidence, err := strconv.ParseFloat(os.Args[2], 64)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
	percentageOfPointsOnPlane, err := strconv.ParseFloat(os.Args[3], 64)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
	numberOfIterations := GetNumberOfIterations(confidence, percentageOfPointsOnPlane)
	nb := numberOfIterations 

	// Read eps from 4th argument from user.
	eps, err := strconv.ParseFloat(os.Args[4], 64)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

	// Create and start the RANSAC find dominant plane pipeline. This pipeline automatically stops
	// after the required number of iterations.
	rand.Seed(time.Now().UnixNano())
	debut := time.Now(); // chrono
	fmt.Println("Debut")

	stop := make(chan bool)
	defer close(stop)

	var wg sync.WaitGroup

	wg.Add(4)
	plane := planeEstimator(&wg, stop , takeN(&wg, stop, triplet(&wg, stop, randomGenerator(&wg, stop, randomNumbers, originalPoints), nb), numberOfIterations), GetPlane)

	fanOut:= 3
	wg.Add(fanOut)
	supportPlane:= make([]<-chan Plane3DwSupport, fanOut)
	for i:=0; i<fanOut; i++ {
		supportPlane[i]= supportingPointFinder(&wg, stop, plane, GetSupport, eps, originalPoints)
	}

	wg.Add(2)
	for i := range DominantPlanIdentifier(&wg, stop, fanIn(&wg,stop,supportPlane), &bestSupport){
		fmt.Println(i)
	}

	stop <- true // stop the thread
	wg.Wait() // synchronisation

	fmt.Println("\nFin\n")
	fin := time.Now();
	//fmt.Println(runtime.NumCPU())
	fmt.Printf("Temps d'execution: %s \n", fin.Sub(debut))

	// Once the pipeline has terminated, save the supporting points of the identified dominant plane in
	// a file named by appending _p to the input filename.
	var supportingPoints []Point3D
	supportingPoints = GetSupportingPoints(bestSupport.Plane3D, originalPoints, eps)
    // Extract base filename
    basename := filepath.Base(os.Args[1])
    // Remove extension
    nameWithoutExt := strings.TrimSuffix(basename, filepath.Ext(basename))
	dominantOutputFilenameP1 := nameWithoutExt+"_p1.xyz"
	dominantOutputFilenameP2 := nameWithoutExt+"_p2.xyz"
	dominantOutputFilenameP3 := nameWithoutExt+"_p3.xyz"
	// Check if file exists ("Do not forget to rename the ... _p file to _p1, _p2 et _p3")
	if _, err := os.Stat(dominantOutputFilenameP1); os.IsNotExist(err) {
		// File _p1 does not already exist
		SaveXYZ(dominantOutputFilenameP1, supportingPoints)
	} else if _, err := os.Stat(dominantOutputFilenameP2); os.IsNotExist(err) {
		// File _p2 does not already exist
		SaveXYZ(dominantOutputFilenameP2, supportingPoints)
	} else if _, err := os.Stat(dominantOutputFilenameP3); os.IsNotExist(err) {
		// File _p3 does not already exist
		SaveXYZ(dominantOutputFilenameP3, supportingPoints)
	}

	// Save the original point cloud without the supporting points of the dominant plane in a file named
	// by appending _p0 to the input filename.
	cloudWithoutsupportPoints := RemovePlane(bestSupport.Plane3D, originalPoints, eps)
	cloudWithoutsupportPointsFilename := nameWithoutExt+"_p0.xyz"
	SaveXYZ(cloudWithoutsupportPointsFilename, cloudWithoutsupportPoints)

}
