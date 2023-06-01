% Name: Anouar Smaili

% predicate that reads the point cloud in a file and creates a list of 3D points
read_xyz_file(File, Points) :-
    open(File, read, Stream),
    read_xyz_points(Stream,Points),
    close(Stream).

read_xyz_points(Stream, []) :-
    at_end_of_stream(Stream).
read_xyz_points(Stream, [Point|Points]) :-
    \+ at_end_of_stream(Stream),
    read_line_to_string(Stream,L), 
    split_string(L, "\t", "\s\t\n",XYZ), 
    convert_to_float(XYZ,Point),
    read_xyz_points(Stream, Points).

convert_to_float([],[]).
convert_to_float([H|T],[HH|TT]) :-
    atom_number(H, HH),
    convert_to_float(T,TT).


% This predicate should be true if Point3 is a triplet of points randomly selected from
% the list of points Points. The triplet of points is of the form \[[x1,y1,z1], [x2,y2,z2], [x3,y3,z3]].
random3points(Points, Point3) :-
    length(Points, Len),
    Len >= 3,
    random_permutation(Points, RandomPoints),
    RandomPoints = [[X1,Y1,Z1], [X2,Y2,Z2], [X3,Y3,Z3] | _],
    Point3 = [[X1,Y1,Z1], [X2,Y2,Z2], [X3,Y3,Z3]].

% Cas de tests pour le predicat plane
test(random3points, 1) :- 
    random3points([[1,2,3],[4,5,6],[7,8,9],[10,11,12]], Point3),
    member([1,2,3], Point3);member([4,5,6], Point3);member([7,8,9], Point3);member([10,11,12], Point3).
test(random3points, 2) :- \+random3points([[1,2,3],[4,5,6]], false). % there are not enough points in Points to select three random points.
test(random3points, 3) :- \+random3points([], false) % there are no points in Points


% This predicate should be true if Plane is the equation of the plane defined by the three
% points of the list Point3. The plane is specified by the list \[a,b,c,d] from the
% equation ax+by+cz=d. The list of points is of the form [[x1,y1,z1], [x2,y2,z2], [x3,y3,z3]].
plane(Point3 , Plane) :- 
    Point3 = [[X1,Y1,Z1], [X2,Y2,Z2], [X3,Y3,Z3]],
    A1 is X2 - X1,
    B1 is Y2 - Y1,
    C1 is Z2 - Z1,
    A2 is X3 - X1,
    B2 is Y3 - Y1,
    C2 is Z3 - Z1,
    A is (B1*C2) - (B2*C1),
    B is (A2*C1) - (A1*C2),
    C is (A1*B2) - (B1*A2),
    D is -A*X1-B*Y1-C*Z1,
    Plane = [A,B,C,D].

% Cas de tests pour le predicat plane
test(plane, 1) :- plane([1,5,3], [4,5,2], [7,1,9], [-4, -24, -12, 160]).
test(plane, 2) :- plane([1,2,-2], [3,-2,1], [5,1,-4], [11, 16, 14, -15]).
test(plane, 3) :- plane([1,2,3], [3,3,3], [4,4,5], [2, -4, 1, 3]).

% This predicate should be true if the support of plane Plane is composed of N points
% from the list of points Point3 when the distance Eps is used.
support(Plane, Points, Eps, N) :- support(Plane, Points, Eps, N, 0).
support(_,[], _, N, N).
support(Plane, [[X1,Y1,Z1]|L], Eps, N, T) :-
    Plane = [A, B, C, D],
    Distance is abs(A*X1 + B*Y1 + C*Z1 + D) / sqrt(A*A + B*B + C*C),
    Distance =< Eps,
    TT is T+1, !,
    support(Plane, L, Eps, N, TT); support(Plane, L, Eps, N, T).

% Cas de tests pour le predicat support
test(support, 1) :- support([3,2,1,5],[[2,2,2],[425,442,796],[1,1,4],[1,1,1]], 10, 3).
test(support, 2) :- support([3,2,1,5],[[2,2,2],[425,442,796],[1,1,4],[1,1,1]], 1000, 4).
test(support, 3) :- support([3,2,1,5],[[2,2,2],[425,442,796],[1,1,4],[1,1,1]], 1000, 0).


% This predicate should be true if N is the number of iterations required by RANSAC with
% parameters Confidence et Percentage according to the formula given in the
% problem description section.
ransac_number_of_iterations(Confidence, Percentage, N) :-
    N is integer(log(1 - Confidence)/log(1 - Percentage ** 3)).

% Cas de tests pour le predicat ransac_number_of_iterations
test(ransac_number_of_iterations, 1) :- ransac_number_of_iterations(0.99, 0.2, 573).
test(ransac_number_of_iterations, 2) :- ransac_number_of_iterations(0.99, 0.5, 34).
test(ransac_number_of_iterations, 3) :- ransac_number_of_iterations(0.80, 0.1, 1609).

