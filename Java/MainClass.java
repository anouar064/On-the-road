// Name: Anouar Smaili

public class MainClass {

    public static void main(String[] args) {

        PointCloud pointCloud1 = new PointCloud("PointCloud1.xyz");
        PointCloud pointCloud2 = new PointCloud("PointCloud2.xyz");
        PointCloud pointCloud3 = new PointCloud("PointCloud3.xyz");

        // Point Cloud 1
        PlaneRANSAC planeRANSAC1 = new PlaneRANSAC(pointCloud1);
        planeRANSAC1.setEps(0.5);
        int interations1 = planeRANSAC1.getNumberOfIterations(0.99, 5);
        // call the algorithm three times to get three dominant plans
        planeRANSAC1.run(interations1, "PointCloud1_p1.xyz");
        planeRANSAC1.run(interations1, "PointCloud1_p2.xyz");
        planeRANSAC1.run(interations1, "PointCloud1_p3.xyz");
        // Get the PointCloud that contains the points remaining
        PointCloud pointsRemaining1 = planeRANSAC1.getPc();
        pointsRemaining1.save("PointCloud1_p0.xyz");

        // Point Cloud 2
        PlaneRANSAC planeRANSAC2 = new PlaneRANSAC(pointCloud2);
        planeRANSAC2.setEps(0.5);
        int interations2 = planeRANSAC1.getNumberOfIterations(0.99, 5);
        // call the algorithm three times to get three dominant plans
        planeRANSAC2.run(interations2, "PointCloud2_p1.xyz");
        planeRANSAC2.run(interations2, "PointCloud2_p2.xyz");
        planeRANSAC2.run(interations2, "PointCloud2_p3.xyz");
        // Get the PointCloud that contains the points remaining
        PointCloud pointsRemaining2 = planeRANSAC2.getPc();
        pointsRemaining2.save("PointCloud2_p0.xyz");

        // Point Cloud 3
        PlaneRANSAC planeRANSAC3 = new PlaneRANSAC(pointCloud3);
        planeRANSAC3.setEps(0.5);
        int interations3 = planeRANSAC1.getNumberOfIterations(0.99, 5);
        // call the algorithm three times to get three dominant plans
        planeRANSAC3.run(interations3, "PointCloud3_p1.xyz");
        planeRANSAC3.run(interations3, "PointCloud3_p2.xyz");
        planeRANSAC3.run(interations3, "PointCloud3_p3.xyz");
        // Get the PointCloud that contains the points remaining
        PointCloud pointsRemaining3 = planeRANSAC3.getPc();
        pointsRemaining3.save("PointCloud3_p0.xyz");
    }

}
