// Name: Anouar Smaili

import java.util.Iterator;
import java.util.ArrayList;
import java.io.File;
import java.io.FileNotFoundException;
import java.io.PrintWriter;
import java.util.Random;
import java.util.Scanner;

public class PointCloud {

    // Store the points belonging to the Cloud
    private ArrayList<Point3D> points;

    // A constructor from a xyz file
    public PointCloud(String filename) {
        points = new ArrayList<Point3D>();
        try {
            Scanner sc = new Scanner(new File(filename));
            // Skip "x y z"
            sc.next();
            sc.next();
            sc.next();
            while (sc.hasNext()) {
                double x = sc.nextDouble();
                double y = sc.nextDouble();
                double z = sc.nextDouble();
                Point3D point = new Point3D(x, y, z);
                // Add the points to the cloud
                points.add(point);
            }
            sc.close();
        } catch (FileNotFoundException e) {
            System.out.println("File not found");
        }
    }

    // An empty constructor that constructs an empty point cloud
    public PointCloud() {
        points = new ArrayList<Point3D>();
    }

    // A addPoint method that adds a point to the point cloud
    public void addPoint(Point3D pt) {
        points.add(pt);
    }

    // A getPoint method that returns a random point from the cloud
    public Point3D getPoint() {
        Random rand = new Random();
        int index = rand.nextInt(points.size());
        return points.get(index);
    }

    // A save method that saves the point cloud into a xyz file
    public void save(String filename) {
        try {
            PrintWriter writer = new PrintWriter(new File(filename));
            writer.println("x   y   z");
            for (Point3D point : points) {
                writer.println(point.getX() + "   " + point.getY() + "   " + point.getZ());
            }
            writer.close();
        } catch (FileNotFoundException e) {
            System.out.println("File not found");
        }
    }

    // An iterator method that returns an iterator to the points in the cloud
    public Iterator<Point3D> iterator() {
        // This iterator includes hasNext, next and remove methods
        return points.iterator();
    }
}
