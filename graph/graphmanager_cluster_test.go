/*
 FishDB
*/

package graph

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/Fisch-Labs/FishDB/cluster"
	"github.com/Fisch-Labs/FishDB/cluster/manager"
	"github.com/Fisch-Labs/FishDB/graph/data"
	"github.com/Fisch-Labs/FishDB/graph/graphstorage"
	"github.com/Fisch-Labs/FishDB/hash"
)

func TestClusterWithPhysicalStorage(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	// Define directory paths for cleanup
	dbDir1 := GraphManagerTestDBDir5
	dbDir2 := GraphManagerTestDBDir6

	// 1. Added Automated Cleanup: This block ensures test directories are always removed when the test is done.
	t.Cleanup(func() {
		os.RemoveAll(dbDir1)
		os.RemoveAll(dbDir2)
	})

	dgs1, err := graphstorage.NewDiskGraphStorage(dbDir1, false)
	if err != nil {
		t.Fatalf("Failed to create disk storage for node 1: %v", err)
	}

	ds1, _ := cluster.NewDistributedStorage(dgs1, map[string]interface{}{
		manager.ConfigRPC:           fmt.Sprintf("localhost:%v", 9021),
		manager.ConfigMemberName:    fmt.Sprintf("TestClusterMember-1"),
		manager.ConfigClusterSecret: "test123",
	}, manager.NewMemStateInfo())

	ds1.Start()
	defer ds1.Close()

	dgs2, err := graphstorage.NewDiskGraphStorage(dbDir2, false)
	if err != nil {
		t.Fatalf("Failed to create disk storage for node 2: %v", err)
	}

	ds2, _ := cluster.NewDistributedStorage(dgs2, map[string]interface{}{
		manager.ConfigRPC:           fmt.Sprintf("localhost:%v", 9022),
		manager.ConfigMemberName:    fmt.Sprintf("TestClusterMember-2"),
		manager.ConfigClusterSecret: "test123",
	}, manager.NewMemStateInfo())

	ds2.Start()
	defer ds2.Close()

	err = ds2.MemberManager.JoinCluster(ds1.MemberManager.Name(),
		ds1.MemberManager.NetAddr())
	if err != nil {
		t.Fatalf("Node 2 failed to join cluster with Node 1: %v", err)
	}

	// 2. Reduced Code Duplication: All repetitive test logic is now in the helper function below.
	runClusterReplicationTests(t, ds1, ds2)
}

func TestClusterStorage(t *testing.T) {
	clusterNodes := createCluster(2)
	joinCluster(clusterNodes, t)

	// 2. Reduced Code Duplication: This test now also calls the same helper.
	runClusterReplicationTests(t, clusterNodes[0], clusterNodes[1])
}

/*
Create a cluster with n members (all storage is in memory)
*/
func createCluster(n int) []*cluster.DistributedStorage {
	log.SetOutput(ioutil.Discard)

	var mgs []*graphstorage.MemoryGraphStorage
	var cs []*cluster.DistributedStorage

	cluster.ClearMSMap()

	for i := 0; i < n; i++ {
		mgs = append(mgs, graphstorage.NewMemoryGraphStorage(fmt.Sprintf("mgs%v", i+1)).(*graphstorage.MemoryGraphStorage))
	}

	for i := 0; i < n; i++ {
		ds, _ := cluster.NewDistributedStorage(mgs[i], map[string]interface{}{
			manager.ConfigRPC:           fmt.Sprintf("localhost:%v", 9020+i),
			manager.ConfigMemberName:    fmt.Sprintf("TestClusterMember-%v", i),
			manager.ConfigClusterSecret: "test123",
		}, manager.NewMemStateInfo())
		cs = append(cs, ds)
	}

	return cs
}

/*
joinCluster joins up a given cluster.
*/
func joinCluster(cluster []*cluster.DistributedStorage, t *testing.T) {
	// 4. Marked as Test Helper: Improves failure reporting.
	t.Helper()

	for i, dd := range cluster {
		dd.Start()
		defer dd.Close()

		if i > 0 {
			bootstrapNode := cluster[0]
			err := dd.MemberManager.JoinCluster(bootstrapNode.MemberManager.Name(),
				bootstrapNode.MemberManager.NetAddr())
			if err != nil {
				// 3. Made Error Messages Descriptive
				t.Fatalf("member %q failed to join cluster with %q: %v",
					dd.MemberManager.Name(), bootstrapNode.MemberManager.Name(), err)
			}
		}
	}
}

// runClusterReplicationTests contains all the shared logic for verifying
// data replication across different layers of the database.
func runClusterReplicationTests(t *testing.T, ds1, ds2 *cluster.DistributedStorage) {
	// 4. Marked as Test Helper: Improves failure reporting.
	t.Helper()

	// *** Direct storage
	sm1 := ds1.StorageManager("foo", true)
	sm2 := ds2.StorageManager("foo", true)

	loc, err := sm1.Insert("test123")
	if loc != 1 || err != nil {
		t.Fatalf("Direct storage insert failed: got loc=%d, err=%v, want loc=1, err=nil", loc, err)
	}
	cluster.WaitForTransfer()
	var res string
	if err := sm2.Fetch(1, &res); err != nil || res != "test123" {
		t.Fatalf("Direct storage fetch failed: got %q, err=%v, want 'test123'", res, err)
	}

	// *** HTree storage
	sm1 = ds1.StorageManager("foo2", true)
	sm2 = ds2.StorageManager("foo2", true)

	htree, err := hash.NewHTree(sm1)
	if err != nil {
		t.Fatalf("Failed to create HTree: %v", err)
	}
	if _, err := htree.Put([]byte("123"), "TestValue"); err != nil {
		t.Fatalf("HTree put failed: %v", err)
	}
	cluster.WaitForTransfer()
	htree2, _ := hash.LoadHTree(sm2, 1)
	if val, err := htree2.Get([]byte("123")); err != nil || val != "TestValue" {
		t.Fatalf("HTree get failed: got %q, err=%v, want 'TestValue'", val, err)
	}

	// *** GraphManager storage
	gm1 := NewGraphManager(ds1)
	gm2 := NewGraphManager(ds2)

	nodeMap := map[string]interface{}{"key": "123", "kind": "testnode", "foo": "bar"}
	if err := gm1.StoreNode("main", data.NewGraphNodeFromMap(nodeMap)); err != nil {
		t.Fatalf("GraphManager store node failed: %v", err)
	}
	cluster.WaitForTransfer()
	node, err := gm2.FetchNode("main", "123", "testnode")
	if err != nil || node.Attr("foo") != "bar" {
		t.Fatalf("GraphManager fetch node failed: node=%v, err=%v", node, err)
	}
}
