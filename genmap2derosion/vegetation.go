package genmap2derosion

/*
func (w *World) growTrees() {
  // Random Position
  {
    i := rand.Int()%(w.dim.x*w.dim.y)
    n := surfaceNormal(i, w.heightmap, w.dim, w.scale)

    if w.waterpool[i] == 0.0 && w.waterpath[i] < 0.2 && n.y > 0.8 {

      Plant ntree(i, w.dim)
      ntree.root(w.plantdensity, w.dim, 1.0)
      w.trees.push_back(ntree)
    }
  }

  // Loop over all Trees
  for i := 0; i < w.trees.size(); i++ {

    // Grow the Tree
    w.trees[i].grow()

    // Spawn a new Tree!
    if rand()%50 == 0 {
      // Find New Position
      glm::vec2 npos = w.trees[i].pos + glm::vec2(rand()%9-4, rand()%9-4)

      // Check for Out-Of-Bounds
      if npos.x >= 0 && npos.x < dim.x && npos.y >= 0 && npos.y < dim.y {

        Plant ntree(npos, w.dim)
        n := surfaceNormal(ntree.index, w.heightmap, w.dim, w.scale)

        if w.waterpool[ntree.index] == 0.0 && w.waterpath[ntree.index] < 0.2 && n.y > 0.8 && (double)(rand()%1000)/1000.0 > w.plantdensity[ntree.index] {
          ntree.root(w.plantdensity, w.dim, 1.0)
          w.trees.push_back(ntree)
        }
      }
    }

    // If the tree is in a pool or in a stream, kill it.
    if w.waterpool[w.trees[i].index] > 0.0 || w.waterpath[w.trees[i].index] > 0.2 || rand()%1000 == 0 {
      //Random Death Chance
      w.trees[i].root(w.plantdensity, w.dim, -1.0)
      w.trees.erase(w.trees.begin()+i)
      i--
    }
  }
}

type Plant struct {
  pos vec2
  index int
  size float32
  maxsize float32
  rate float32

  Plant& operator=(const Plant& o){
    if this != &o {  //Self Check
      pos = o.pos
      index = o.index
      size = o.size
    }
    return *this
  }
}

func (p *Plant) init() {
  p.size = 0.5
  p.maxsize = 1.0
  p.rate = 0.05
}

func newPlant(i int, d ivec2) Plant{
  var p Plant
  p.index = i
  p.pos = vec2(i/d.y, i%d.y)
  p.init()
  return p
}

func newPlantt(p vec2,  d ivec2) Plant{
  var p Plant
  p.pos = p
  p.index = int(p.x)*d.y+int(p.y)
  p.init()
  return p
}


func (p *Plant) grow() {
  p.size += p.rate*(p.maxsize-p.size)
}

func (p *Plant) root(density []float64, dim ivec2, f float64) {

  //Can always do this one
  density[index] += f*1.0

  if pos.x > 0 {
    //
    density[index - 256] += f*0.6      //(-1, 0)

    if pos.y > 0 {
      density[index - 257] += f*0.4    //(-1, -1)
    }
    if pos.y < 256-1 {
      density[index - 255] += f*0.4    //(-1, 1)
    }
  }

  if pos.x < 256-1 {
    //
    density[index + 256] += f*0.6    //(1, 0)

    if pos.y > 0 {
      density[index + 255] += f*0.4    //(1, -1)
    }
    if pos.y < 256-1 {
      density[index + 257] += f*0.4    //(1, 1)
    }
  }

  if pos.y > 0 {
    density[index - 1]   += f*0.6    //(0, -1)
  }
  if pos.y < 256-1 {
    density[index + 1]   += f*0.6    //(0, 1)
  }
}*/
