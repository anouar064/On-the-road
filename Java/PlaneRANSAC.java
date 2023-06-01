// Name: Anouar Smaili

import java.util.Iterator;
import java.lang.Double;

public class PlaneRANSAC {

    // pour le pourcentage 5% est bien
    private double eps;
    private int bestSupport;
    private PointCloud pc;

    // A constructor that takes as input a point cloud
    public PlaneRANSAC(PointCloud pc) {
        this.pc = pc;
    }

    // Setter for the epsilon value
    public void setEps(double eps) {
        this.eps = eps;
    }

    // Getter for the epsilon value
    public double getEps() {
        return eps;
    }

    // Get the point cloud
    public PointCloud getPc() {
        return pc;
    }

    // Returns the estimated number of iterations required to obtain
    // a certain level of confidence to identify a plane made of
    // a certain percentage of points
    public int getNumberOfIterations(double confidence,
            double percentageOfPointsOnPlane) {

        return (int) ((Math.log(1 - confidence)) / (Math.log(1 - Math.pow(percentageOfPointsOnPlane, 3))));
    }

    // Runs the RANSAC algorithm for identifying the dominant plane of the
    // point cloud (only one plane)
    public void run(int numberOfIterations,
            String filename) {
        // Filename being the xyz file that will contain the points of the dominant
        // plane.
        // This method will also remove the plane points from the point cloud

        // PointCloud pcToRemove that contains the points belonging to the dominant
        // plane
        PointCloud pcToRemove = new PointCloud();

        // Initially, no dominant plane has been found, and the best support is set
        // to 0
        bestSupport = 0;

        // Repeat n times so that we are confident to have found the dominant plane.
        for (int i = 0; i <= numberOfIterations; i++) {

            // PointCloud pcPossibleToRemove that contains the points belonging to a
            // potential dominant plane
            PointCloud pcPossibleToRemove = new PointCloud();

            // Randomly draw 3 points from the point cloud.
            Point3D point1 = pc.getPoint();
            Point3D point2 = pc.getPoint();
            Point3D point3 = pc.getPoint();

            // Compute the plane equation from these 3 points.
            Plane3D plane = new Plane3D(point1, point2, point3);

            // Count the number of points that are at a distance less than eps (ε) from
            // that plane. This number is the support for the current plane.
            Iterator<Point3D> it = pc.iterator();
            int support = 0;
            while (it.hasNext()) {
                Point3D point = it.next();
                double distance = plane.getDistance(point);
                if (distance < getEps()) {
                    pcPossibleToRemove.addPoint(point);
                    support++;
                }
            }

            // If the current support is higher than the best support value, then the
            // current plane becomes the current dominant plane and its support is the new
            // best support.
            if (support > bestSupport) {
                bestSupport = support;
                // dominantPlane = plane
                pcToRemove = pcPossibleToRemove;
            }
        }

        // Remove the points that belong to the dominant plane from the point cloud.
        // Save these points in a new file.
        // Iterator for initial PointCloud
        Iterator<Point3D> initial = pc.iterator();
        // Iterator for the PointCloud that contains points to remove
        Iterator<Point3D> toRemove = pcToRemove.iterator();
        // Take everyPoint from PointCloud toRemove and compare it
        // with every point from the initial PointCloud in order to
        // find it and remove it from the PointCloud
        while (toRemove.hasNext()) {
            Point3D pointToRemove = toRemove.next();
            while (initial.hasNext()) {
                Point3D pointInitial = initial.next();
                // If two points have the same X, Y and Z, then it's the same point
                if ((Double.compare(pointInitial.getX(), pointToRemove.getX()) == 0)
                        && (Double.compare(pointInitial.getY(), pointToRemove.getY()) == 0)
                        && (Double.compare(pointInitial.getZ(), pointToRemove.getZ()) == 0)) {

                    initial.remove();
                }
            }
        }
        // Output the points belonging to the dominant plane in a new file
        pcToRemove.save(filename);

    }
}
