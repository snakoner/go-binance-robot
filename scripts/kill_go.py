import os

temp_file = 'rmfile'

def main():
    os.system('ps aux | grep main.go > %s' % temp_file)
    procs = []
    with open(temp_file, 'r') as f:
        procs = f.read().splitlines()

    try:
        os.remove(temp_file)
    except FileNotFoundError:
        print('Cant remove file %s' % temp_file)

    procs_id = []
    for p in procs:
        s = p.split(' ')
        s = [x for x in s if x]
        if 'grep' not in s:
            procs_id.append(s[1])
    
    for pid in procs_id:
        os.system('kill -9 %s' % pid)

    print('Killed %d processes' % len(procs_id))

if __name__ == '__main__':
    main()