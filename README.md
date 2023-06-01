# On the road Project

The intelligent vehicles of the future will be equipped with a multitude of sensors in order to capture information about the surrounding scene and thus being able to autonomously navigate. One of these sensors is the Laser Scanner or LiDAR (Light Detection And Ranging). Using a LiDAR, the vehicle can scan the scene in front by sweeping few laser beams (typically between 8 to 64 lasers).

Each time a laser beam hit an object, the laser light bounces back to the LiDAR from which a precise distance can be estimated. A complete scan of the scene with these lasers will therefore generate a large set of 3D points (also called point cloud) that correspond to the structure of the scene. The figure below shows a typical point cloud captured by a car equipped with a LiDAR. A view of the same scene captured by a camera is also shown.

In this comprehensive assignment, we need to implement an algorithm that will detect the dominant planes in a cloud of 3D points. The dominant plane of a set of 3D points is the plane that contains the largest number of points. A 3D point is contained in a plane if it is located at a distance less than eps (ε) from that plane. To detect these planes, we will use the RANSAC algorithm. More specifically, three point clouds will be given to us. For each of them we need to find the three most dominant planes.

## RANSAC : Random Sampling Consensus:

RANSAC is an iterative algorithm that is used to identify a geometric entity (or model) from a set of data that contains a large amount of outliers (data that do not belong to the model). It proceeds by randomly drawing the minimum number of samples required to estimate the parameters of a model instance and then validate it by counting the number of additional samples that support the computed model. In our case, we are looking for a planar structure, made of several points, while the majority of the points in the set are outside that plane. The seek geometric entity is therefore a plane of the form: ax+by+cz=d. A minimum of 3 points is required to compute the equation of a plane.

RANSAC for the case of plane identification in 3D proceeds as follows :

1. Initially, no dominant plane has been found, and the best support is set to 0 (see Step 5.)
2. Randomly draw 3 points from the point cloud.
3. Compute the plane equation from these 3 points.
4. Count the number of points that are at a distance less than eps (ε) from that plane. This number
   is the support for the current plane.
5. If the current support is higher than the best support value, then the current plane becomes the
   current dominant plane and its support is the new best support.
6. Repeat 2 to 5 until we are confident to have found the dominant plane.
7. Remove the points that belong to the dominant plane from the point cloud. Save these points in a
   new file.
   Step 6. raises the following question: how many iterations should we perform if we want to be almost
   certain (let’s say at 99%) that we have found the dominant plane?
   First, suppose that the percentage of points that support the dominant plane is p% of the total number of points in the cloud. The probability of randomly picking three points that belong to this plane is therefore p^3 %. We can then conclude that the probability of picking a set of random that contains at least one outlier is (1- p^3 )%. If we pick k random triplets of points, the probability that these sets always contains an outlier is (1- p^3 )^k %. Consequently, the probability of finding at least one set made of 3 points that belongs to the dominant plane is 1-(1- p^3)^k %. We must therefore find the value of k that give us a confidence probability of, let’s say, C= 99%.
   k= log( 1 - C ) / log( 1- p^3 )
   To find the three most dominant plane, we then need to repeat the complete procedure 3 times.
