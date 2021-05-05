import os
import sys

if len(sys.argv) != 6:
  sys.stderr.write('Arguments error. Usage:\n')
  sys.stderr.write('\tpython evaluate.py model features scores prc roc\n')
  sys.exit(1)

model_file = sys.argv[1]
matrix_file = os.path.join(sys.argv[2], 'test.pkl')
scores_file = sys.argv[3]
prc_file = sys.argv[4]
roc_file = sys.argv[5]

print(model_file, matrix_file, scores_file, prc_file, roc_file)
