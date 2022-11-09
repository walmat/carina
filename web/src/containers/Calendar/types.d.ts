interface ReleasesRoot {
  [date: string]: Releases;
}

interface Releases {
  [id: string]: Release;
}

interface Release {
  name: string;
  image: string;
  stores: Stores;
}

interface Stores {
  [id: string]: ReleaseEntry;
}

interface ReleaseEntry {
  name: string;
  date: string;
  datum: Datum;
}

type Datum = Data[];

interface Data {
  key: string;
  value: string;
}
