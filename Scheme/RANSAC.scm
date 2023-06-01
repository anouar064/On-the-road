; Name: Anouar Smaili

#lang scheme
; function that reads the point cloud in a file and creates a list of 3D points
(define (readXYZ fileIn)
(let ((sL (map (lambda s (string-split (car s)))
(cdr (file->lines fileIn)))))
(map (lambda (L)
(map (lambda (s)
(if (eqv? (string->number s) #f)
s
(string->number s))) L)) sL)))


; function that is to be able to pick 3 random points from a list of points.
; input: list of points Ps
; output: 3 random points p1 p2 p3
(define (random-3-points Ps)
  (let ((len (length Ps)))
    (if (< len 3)
        (error "La liste doit contenir au moins 3 points")
        (let* ((p1 (list-ref Ps (random len)))
               (p2 (list-ref Ps (random len)))
               (p3 (list-ref Ps (random len))))
          (list p1 p2 p3)))))

; Function that computes a plane equation from 3 points. 
; The equation is described by a list of the four parameters, 
; from the equation ax+by+cz=d
; input: 3 points P1 P2 P3
; output: the four parameters a b c d from the plane equation 
(define (plane P1 P2 P3)
  (let* ((x1 (car P1))
         (y1 (cadr P1))
         (z1 (caddr P1))
         (x2 (car P2))
         (y2 (cadr P2))
         (z2 (caddr P2))
         (x3 (car P3))
         (y3 (cadr P3))
         (z3 (caddr P3))
         (a1 (- x2 x1)) (b1 (- y2 y1)) (c1 (- z2 z1))
         (a2 (- x3 x1)) (b2 (- y3 y1)) (c2 (- z3 z1))
         (a (- (* b1 c2) (* b2 c1)))
         (b (- (* a2 c1) (* a1 c2)))
         (c (- (* a1 b2) (* b1 a2)))
         (d (- (* (* a x1) -1) (* b y1) (* c z1))))
    (list a b c d)))


; function that returns the distance from a plane to a point
; input: a plane and a point
; output: the distance between them
(define (distance-from-plane plane point)
  (let ((a (list-ref plane 0))
        (b (list-ref plane 1))
        (c (list-ref plane 2))
        (d (list-ref plane 3))
        (x (list-ref point 0))
        (y (list-ref point 1))
        (z (list-ref point 2)))
    (/ (abs (+ (* a x) (* b y) (* c z) d))
       (sqrt (+ (* a a) (* b b) (* c c))))))    

; function that returns the support count and the plane parameter in a pair
; this is the helper function doing the work
; input: plane, point, eps, count (incremented each time it supports the plane)
; output: count which is the support-size
(define (support-helper plane points eps count)
  (cond ((null? points) count)
        (else
          (if (<= (distance-from-plane plane (car points)) eps)
              (support-helper plane (cdr points) eps (+ count 1))
              (support-helper plane (cdr points) eps count)))))

; function that returns the support count and the plane parameter in a pair
; input: plane, points, eps
; output: the plane parameter and the support count in a pair
(define (support plane points eps)
  (let ((support-size (support-helper plane points eps 0)))
    (cons plane support-size)))

; function that repeat the random sampling K times in order to find the
; dominant plane (the plane with the best support)
; this is the helper function doing the work
; input: list of points Ps, k, eps, bestSupport, count
; output: bestSupport
(define (dominantPlane-helper Ps k eps bestSupport count)
  (cond ((= k count) bestSupport)
        (else
         (let* ((resultRandom (random-3-points Ps))
                (resultPlane (plane (car resultRandom) (cadr resultRandom) (caddr resultRandom)))
                (currentSupport (support resultPlane Ps eps)))
           (if (> (cdr currentSupport) (cdr bestSupport))
               (dominantPlane-helper Ps k eps currentSupport (+ 1 count))
               (dominantPlane-helper Ps k eps bestSupport (+ 1 count)))))))

; function that repeat the random sampling K times in order to find the
; dominant plane (the plane with the best support)
; input: list of points Ps, k, eps
; output: plane
(define (dominantPlane Ps k eps)
    (let ((plane (dominantPlane-helper Ps k eps '((0 0 0 0) . 0) 1)))
      plane))


; function that computes this number based on certain confidence 
; and a certain percentage of points
; input: confidence, percentage
; output: number of iterations
(define (ransacNumberOfIterations confidence percentage)
  (define (log-base-10 x) (/ (log x) (log 10)))
  (define iterations 
    (exact-round 
     (/ (log-base-10 (- 1 confidence))
        (log-base-10 (- 1 (expt percentage 3))))))
  iterations)


; main function
; This function will return a pair made of the dominant plane equation and the number of points that
; supports it. The file should be read in a let statement such that the variable Ps will contain the list of
; 3D points in the rest of the function.
; input: filename, confidence, percentage, eps
; output: the dominant plane equation and the number of points that supports it
(define (planeRANSAC filename confidence percentage eps)
  (let ((Ps (readXYZ filename))
        (k (ransacNumberOfIterations confidence percentage)))
    (dominantPlane Ps k eps)))

; Test cases

(random-3-points '((0 0 0) (1 1 1) (2 2 2) (3 3 3) (4 4 4)))
(random-3-points '((0 0 0) (1 1 1) (2 2 2) (3 3 3) (4 4 4)))

(plane '(1 5 3) '(4 5 6) '(8 1 9))
(plane '(1 2 3) '(6 7 10) '(1 5 6))

(distance-from-plane '(-6 -15 15 -9) '(1 2 3))
(distance-from-plane '(-6 -15 15 -9) '(0 9 5))

(support '(-6 -15 15 -9) '((1 2 3) (6 7 10) (1 5 6) (13 15 119)) 2)
(support '(-6 -15 15 -9) '((1 2 3) (6 7 10) (1 5 6) (13 15 119)) 100)

(dominantPlane '((-5.3129109 0.336153973 0.39271088)
  (-5.442733358 0.149305066 -0.048846397)
  (-5.483027127 0.355853761 -0.398578564)
  (-6.617345499 -0.045837601 0.603870467)
  (-5.248219829 0.351818111 0.003580368)
  (-5.468973216 0.371794361 -0.397638032)
  (-5.644454252 0.184573817 0.568371749)
  (-5.279517288 0.368245198 0.196743798)
  (-5.323459657 0.181146444 -0.241224688)
  (-5.922420603 0.426881619 0.656796377)
  (-5.704477074 -0.002520161 0.311939149)
  (-5.474087067 0.001976501 -0.099946839)
  (-5.793393512 0.217667621 -0.476130388)
  (-5.637250032 0.218972337 0.567770934)
  (-5.220361402 0.396320827 0.194626284)
  (-5.396068797 0.018503336 -0.098522936)
  (-5.805644215 0.23596665 -0.47719445)
  (-5.278297149 0.415341469 0.390577362)
  (-5.406017017 0.033047476 0.097667592)
  (-5.354502073 0.231558902 -0.24271765)
  (-5.477759175 0.246447385 0.551849282)
  (-5.275596747 0.433087031 0.196780947)
  (-5.392854217 0.051583954 -0.098468167)
  (-6.774190261 0.075294669 0.868087728)
  (-5.405589635 0.260577014 0.347314024)
  (-5.18007668 0.443161437 0.003538847)
  (-5.555133572 0.072702503 -0.305451924)
  (-5.597148889 0.081438423 0.510813549)
  (-5.359771214 0.277666415 0.148037866)
  (-5.200039862 0.463139249 -0.187782099)
  (-5.493142283 0.297822474 0.553651726)
  (-5.300247315 0.484260767 0.197858286)
  (-5.453130253 0.102360923 -0.099581731)
  (-5.52072789 0.309272786 -0.454112313)
  (-5.297120437 0.498687866 0.392490164)
  (-5.306073501 0.307540385 0.146604091)
  (-5.499771407 0.120116483 -0.100439725))
                      5 1)

(ransacNumberOfIterations 0.99 0.2)
(ransacNumberOfIterations 0.8 0.5)

(planeRANSAC "Point_Cloud_1_No_Road_Reduced.xyz" 0.99 0.2 1)
(planeRANSAC "Point_Cloud_1_No_Road_Reduced.xyz" 0.99 0.5 2)
(planeRANSAC "Point_Cloud_1_No_Road_Reduced.xyz" 0.80 0.2 0.5)

(planeRANSAC "Point_Cloud_2_No_Road_Reduced.xyz" 0.99 0.2 1)
(planeRANSAC "Point_Cloud_2_No_Road_Reduced.xyz" 0.99 0.5 2)
(planeRANSAC "Point_Cloud_2_No_Road_Reduced.xyz" 0.80 0.2 0.5)

(planeRANSAC "Point_Cloud_3_No_Road_Reduced.xyz" 0.99 0.2 1)
(planeRANSAC "Point_Cloud_3_No_Road_Reduced.xyz" 0.99 0.5 2)
(planeRANSAC "Point_Cloud_3_No_Road_Reduced.xyz" 0.80 0.2 0.5)