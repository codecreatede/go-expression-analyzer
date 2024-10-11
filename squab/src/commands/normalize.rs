use std::{
    fs::File,
    io::{self, BufReader},
    path::{Path, PathBuf},
};

use thiserror::Error;
use tracing::info;

use crate::{
    count,
    normalization::{self, calculate_fpkms, calculate_tpms},
    read_features,
};

#[derive(Debug, Error)]
pub enum NormalizeError {
    #[error("I/O error")]
    Io(#[source] io::Error),
    #[error("could not open file: {1}")]
    OpenFile(#[source] io::Error, PathBuf),
    #[error("invalid counts")]
    ReadCounts(#[source] count::reader::ReadCountsError),
    #[error("invalid annotations")]
    ReadAnnotations(#[source] crate::ReadFeaturesError),
    #[error("normalization error")]
    Normalization(#[source] normalization::Error),
}

pub fn normalize<P, Q>(
    counts_src: P,
    annotations_src: Q,
    feature_type: &str,
    id: &str,
    method: normalization::Method,
) -> Result<(), NormalizeError>
where
    P: AsRef<Path>,
    Q: AsRef<Path>,
{
    let counts_src = counts_src.as_ref();
    let annotations_src = annotations_src.as_ref();

    let mut reader = File::open(counts_src)
        .map(BufReader::new)
        .map(count::Reader::new)
        .map_err(|e| NormalizeError::OpenFile(e, counts_src.into()))?;

    let count_map = reader.read_counts().map_err(NormalizeError::ReadCounts)?;

    let mut gff_reader = crate::gff::open(annotations_src)
        .map_err(|e| NormalizeError::OpenFile(e, annotations_src.into()))?;

    let feature_map = read_features(&mut gff_reader, feature_type, id)
        .map_err(NormalizeError::ReadAnnotations)?;

    let feature_ids: Vec<_> = feature_map.keys().map(|id| id.into()).collect();

    let values = match method {
        normalization::Method::Fpkm => {
            info!("calculating fpkms");
            calculate_fpkms(&count_map, &feature_map).map_err(NormalizeError::Normalization)?
        }
        normalization::Method::Tpm => {
            info!("calculating tpms");
            calculate_tpms(&count_map, &feature_map).map_err(NormalizeError::Normalization)?
        }
    };

    let stdout = io::stdout().lock();
    let mut writer = normalization::Writer::new(stdout);

    writer
        .write_values(&feature_ids, &values)
        .map_err(NormalizeError::Io)?;

    Ok(())
}
